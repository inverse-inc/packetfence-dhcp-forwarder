dhcp-forwarder
==============
Part of the PacketFence-DHCP-Forwarder, 'dhcp-forwarder' sends captured UDP packets from a listening device to a destination host and port. Takes it's configuration from the DHCP-Forwarder.toml file.


Microsoft Windows based system
------------------------------
An installer is provided for an easy step-by-step installation, setup and configuration of the PacketFence-DHCP-Forwarder on Microsoft Windows based system.

More information about it can be found [here](https://github.com/inverse-inc/packetfence-dhcp-forwarder#dhcp-forwarder).


Linux based system
------------------
The 'dhcp-forwarder' can work on a Linux based system as well. Unfortunately, there is no step-by-step installer, therefore manual installation is required.

### Golang environment

Developed in Golang, 'dhcp-forwarder' requires architecture based compilation to get a working binary.

Theses instructions assumes you already have a working Golang environment. If it is not the case, instructions about setting such environment can be found [here](https://golang.org/doc/install).

### Build the binary

 * Get the 'dhcp-forwarder' sources (either by forking and cloning the project or by downloading the archive from the [Github repository](https://github.com/inverse-inc/packetfence-dhcp-forwarder))
 * Make sure the 'dhcp-forwarder' working directory (./dhcp-forwarder/) is part of the GOPATH and that the source path is right (GOPATH/src/github.com/inverse-inc/packetfence-dhcp-forwarder/dhcp-forwarder)
 * Within the 'dhcp-forwarder' working directory (./dhcp-forwarder/), fetch required external librairies using :
```
go get ./...
```
 * Within the 'dhcp-forwarder' working directory (./dhcp-forwarder/), build the binary using :
```
go build
```
 * Outputed 'dhcp-forwarder' is the newly built binary file
