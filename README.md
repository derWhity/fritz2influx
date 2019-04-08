# fritz2influx

A simple tool to read the current traffic information from a Fritz!Box and write it to an [InfluxDB](https://www.influxdata.com/) instance.

## Installation

The built version of `fritz2influx` is self-containing. Just download the binary release from the "Releases" section, extract to the directory of your choice and run it.

### Building it yourself

The only prerequisite to build `fritz2influx` is a recent installation of the `go` tools obtainable at https://golang.org/dl. Once this is on your system, clone this repository into a directory of your choice
```
> git clone https://github.com/derWhity/fritz2influx.git
> cd fritz2influx
```
and build it via
```
> go build
```
This will download all dependencies for the project and build a binary for your system.

## Usage

You can directly execute `fritz2influx` from within its directory by running it without parameters. It will then try to load its configuration file from a default location (`/etc/fritz2influx/fritz2influx.conf`) and run the data collection. The Fritz!Box in your network will automatically be discovered via UPnP and the traffic measurements are available without providing username and password - so there is no need to configure this part. The connection to the InfluxDB instance has to be configured in the configuration file, though (the defaults are only good for development and should **not** be used for a final installation).

For setting up a configuration file, see below.

The `fritz2influx` command has two parameters:

* `-c` allows you to select a configuration file to use - e.g. `./fritz2influx -c /home/pi/fritz2influx.conf`
* `-dump` will dump fritz2influx's default configuration to the standard output. You can use it to write a new configuration file:
```
> ./fritz2influx -dump
influx:
  addr: http://localhost:8086
  username: ""
  password: ""
  database: fritzBox
  measurement: WANConnection
collection:
  discoveryInterval: 1h0m0s
  discoveryCooldown: 30s
  interval: 10s
>
```

### Configuration

`fritz2influx` searches for a configuration file at `/etc/fritz2influx/fritz2influx.conf` by default.
If your configuration file is located in this location, simply start `fritz2influx` without any parameters:

```
./fritz2influx
```

In order to create a configuration file, you can obtain the default configuration by using the `-dump` parameter.
This will output the configuration to stdout. This means, you can create a default configuration file by running:

```
./fritz2influx -dump > /etc/fritz2influx/fritz2influx.conf
```

After this, edit it with the text editor of your choice.
