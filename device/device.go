package device

import (
	"fmt"

	"github.com/huin/goupnp/dcps/internetgateway2"
	"github.com/huin/goupnp/soap"
	"github.com/sirupsen/logrus"
)

const (
	// Custom function on a Fritz!Box used here to read the transfer rates
	getAddonInfos = "GetAddonInfos"

	valByteSendRate       = "byteSendRate"
	valByteReceiveRate    = "byteReceiveRate"
	valPacketSendRate     = "packetSendRate"
	valPacketReceiveRate  = "packetRecieveRate"
	valTotalBytesSent     = "totalBytesSent"
	valTotalBytesReceived = "totalBytesReceived"
)

// Discover runs a discovery on the local network for routers providing the service
// "urn:schemas-upnp-org:service:WANCommonInterfaceConfig:1" via UPnP and returns a list
// of all of those devices found
func Discover(logger *logrus.Entry) ([]*Device, error) {
	logger.Info("Discovering routers in the local network...")
	clients, errs, err := internetgateway2.NewWANCommonInterfaceConfig1Clients()
	if err != nil {
		return nil, err
	}
	if len(errs) > 0 {
		for _, err := range errs {
			logger.WithError(err).Errorln("Error occured during device discovert")
		}
	}
	// De-duplicate the devices based on their host name (IP in most cases)
	deviceMap := make(map[string]*internetgateway2.WANCommonInterfaceConfig1)
	out := []*Device{}
	for _, client := range clients {
		host := client.RootDevice.URLBase.Hostname()
		l2 := logger.
			WithField("udn", client.RootDevice.Device.UDN).
			WithField("manufacturer", client.RootDevice.Device.Manufacturer).
			WithField("model", client.RootDevice.Device.ModelName).
			WithField("host", host)
		if _, ok := deviceMap[host]; ok {
			l2.Infof("Ignoring duplicate device at '%s'", host)
		} else {
			l2.Info("New device discovered")
			deviceMap[host] = client
			out = append(out, &Device{
				WANCommonInterfaceConfig1: client,
				Logger:                    l2,
			})
		}
	}
	logger.Infof("Device discovery finished. %d device(s) found", len(out))
	return out, nil
}

// Device represents a router device found during discovery
type Device struct {
	*internetgateway2.WANCommonInterfaceConfig1
	// Logger entry that is preconfigured with fields identifying the router
	Logger *logrus.Entry
}

// FetchReadings loads the transfer readings from the Fritz!Box and unmarshals it into a TransferReadings struct
func (dev *Device) FetchReadings() (*TransferReadings, error) {
	request := interface{}(nil)
	response := &rawReadings{}
	if err := dev.SOAPClient.PerformAction(internetgateway2.URN_WANCommonInterfaceConfig_1, getAddonInfos, request, response); err != nil {
		return nil, err
	}
	return response.toReadings(), nil
}

// TransferReadings represents the network transfer readings requested from the Fritz!Box
type TransferReadings struct {
	ByteSendRate       uint32
	ByteReceiveRate    uint32
	PacketSendRate     uint32
	PacketReceiveRate  uint32
	TotalBytesSent     uint32
	TotalBytesReceived uint32
}

// ToInfluxValues outputs the readings as InfluxDB compatible value map
func (r *TransferReadings) ToInfluxValues() map[string]interface{} {
	return map[string]interface{}{
		valByteSendRate:       r.ByteSendRate,
		valByteReceiveRate:    r.ByteReceiveRate,
		valPacketSendRate:     r.PacketSendRate,
		valPacketReceiveRate:  r.PacketReceiveRate,
		valTotalBytesSent:     r.TotalBytesSent,
		valTotalBytesReceived: r.TotalBytesReceived,
	}
}

func (r *TransferReadings) String() string {
	return fmt.Sprintf("[ Rate: Bytes(⬆ %d / ⬇︎ %d) - Packets(⬆ %d / ⬇︎ %d) | Total: Bytes(⬆ %d / ⬇︎ %d)]",
		r.ByteSendRate,
		r.ByteReceiveRate,
		r.PacketSendRate,
		r.PacketReceiveRate,
		r.TotalBytesSent,
		r.TotalBytesReceived,
	)
}

// rawReadings is the readings structure filled by the SOAP client
type rawReadings struct {
	NewByteSendRate       string
	NewByteReceiveRate    string
	NewPacketSendRate     string
	NewPacketReceiveRate  string
	NewTotalBytesSent     string
	NewTotalBytesReceived string
}

func (r *rawReadings) toReadings() *TransferReadings {
	return &TransferReadings{
		ByteSendRate:       unmashalUint32(r.NewByteSendRate),
		ByteReceiveRate:    unmashalUint32(r.NewByteReceiveRate),
		PacketSendRate:     unmashalUint32(r.NewPacketSendRate),
		PacketReceiveRate:  unmashalUint32(r.NewPacketReceiveRate),
		TotalBytesSent:     unmashalUint32(r.NewTotalBytesSent),
		TotalBytesReceived: unmashalUint32(r.NewTotalBytesReceived),
	}
}

func unmashalUint32(str string) uint32 {
	i, _ := soap.UnmarshalUi4(str)
	return i
}
