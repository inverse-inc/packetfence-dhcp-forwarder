packetfence-forwarder
==============
Part of the PacketFence-Forwarder, 'packetfence-forwarder' sends captured UDP packets from a listening device to a destination host and port. Takes it's configuration from the Packetfence-Forwarder.toml file.


Microsoft Windows based system
------------------------------
An installer is provided for an easy step-by-step installation, setup and configuration of the PacketFence-Packetfence-Forwarder on Microsoft Windows based system.

More information about it can be found [here](https://github.com/inverse-inc/packetfence-forwarder#packetfence-forwarder).


Linux based system
------------------
The 'packetfence-forwarder' can work on a Linux based system as well. Unfortunately, there is no step-by-step installer, therefore manual installation is required.

### Dependencies
 * Requires 'libpcap' as well as 'libpcap-devel'

### Golang environment

Developed in Golang, 'packetfence-forwarder' requires architecture based compilation to get a working binary.

Theses instructions assumes you already have a working Golang environment. If it is not the case, instructions about setting such environment can be found [here](https://golang.org/doc/install).

### Build the binary

 * Get the 'packetfence-forwarder' sources (either by forking and cloning the project or by downloading the archive from the [Github repository](https://github.com/inverse-inc/packetfence--forwarder))
 * Make sure the 'packetfence-forwarder' working directory (./packetfence-forwarder/) is part of the GOPATH and that the source path is right (GOPATH/src/github.com/inverse-inc/packetfence--forwarder/packetfence-forwarder)
 * Within the 'packetfence-forwarder' working directory (./packetfence-forwarder/), fetch required external librairies using :
```
go get ./...
```
 * Within the 'packetfence-forwarder' working directory (./packetfence-forwarder/), build the binary using :
```
go build
```
 * Newly built binary file would be named 'packetfence-forwarder'

### Configuration

 * DestinationHost: To which IP address the trafic is being forwarded (Usually PacketFence management IP address)
 * DestinationPort: To which UDP port of the destination host the trafic is being forwarded (Usually 767)
 * ListeningDevice: Trafic from which network interface should be forwarded (Something like eth0, em1 or any network interface name)
 * Filter: Filter applied to libpcap to forward only certain type of trafic
