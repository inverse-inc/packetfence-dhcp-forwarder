dhcp-forwarder-config-generator:

* Generates DHCP-Forwarder.toml configuration based on user selection of modified "getmac /fo csv /v" output, since it is currently impossible to use gopacket to list human readable interface name.
* Asks for destination ip and port.

The configuration file contains preconfigured BPF for DHCPREQUESTS(3) and DHCPACK(5). That can be modified directly in the configuration file.

Notes: Wireless device selection will not work.

