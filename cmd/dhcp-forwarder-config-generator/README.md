packetfence-forwarder-config-generator
===============================

> Only used for Microsoft Windows based system, Linux based system instructions can be found [here](https://github.com/inverse-inc/packetfence-dhcp-forwarder/tree/master/dhcp-forwarder#linux-based-system)

Part of the PacketFence-Forwarder, 'packetfence-forwarder-config-generator' helps Microsoft Windows based system step-by-step installation.

* Generates Packetfence-Forwarder.toml configuration based on user selection of modified "getmac /fo csv /v" output, since it is currently impossible to use gopacket to list human readable interface name.
* Asks for destination ip and port.

The configuration file contains preconfigured BPF for DHCPREQUESTS(3) and DHCPACK(5). That can be modified directly in the configuration file.

> Note that wireless listening device selection will not work.

