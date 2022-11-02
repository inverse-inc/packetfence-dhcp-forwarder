Packetfence-Forwarder
==============

This tool captures and forwards a subset of DHCP (specifically DHCPREQUEST and DHCPACK) and/or DNS traffic from a Windows x64 DHCP server to a destination IP and port.

Alternatively, for DHCP IP Helpers can be configured on each switch of an infrastructure to forward broadcast only packets.
Those contain all types of DHCP packets but less to none DHCPACK which confirms lease ownership.

DHCP traffic is useful to PacketFence to link MAC adresses and IP addresses, while also helping Fingerbank's fingerprinting process.
Again, the only useful packets to PacketFence are DHCPREQUEST and DHCPACK.

In short, if Packetfence-Forwarder can be deployed centrally, it should be done.
In that case, only useful packets are captured and forwarded from the source, which reduces configuration, transport, processing and storage costs of useless packets, while removing the need to process them to actually neglect them at the destination host.

[Download the installer here.](https://inverse.ca/downloads/PacketFence/windows-packetfence-forwarder/Packetfence-Forwarder-Installer.exe)


Binaries
========
packetfence-forwarder.exe
------------------
This tool captures and forwards DHCP (DHCPREQUEST and DHCPACK, specifically) and/or DNS traffic from a Windows DHCP server to PacketFence.

DHCPREQUEST and DHCPACK packets are the ones being the most important for PacketFence to link MAC to IP addresses and switch location and help Fingerbank to fingerprint the operating system running on the device.

This fingerprinting and localisation process helps a lot in determining violation triggers condition.

With the help of Packetfence-Forwarder, those DHCP packets can be obtained directly and easily from the source.

Alternatively, IP helpers can be configured on each switch to forward DHCP traffic to PacketFence, but only broadcast packets can be captured by them, which is less precise. Deploying Packetfence-Forwarder is simple and centralized.

Packetfence-forwarder is based on gopacket and depends upon WinPCAP to select the requested packets through a BPF, which is really fast. Captured traffic is then forwarded to a configured host and port.

A default BPF is produced by the configuration generator. That filter can be manually modified in the configuration file by the user.

Packetfence-Forwarder requires a Packetfence-Forwarder.toml file to be present in it's working directory (installation directory). Packetfence-Forwarder.toml is generated from packetfence-forwarder-config-generator.exe at installation time, but can be run from the installation directory anytime.

packetfence-forwarder-config-generator.exe
-----------------------------------
Does:

1. Ask for Network Interface Card name and converts it to UUID that will be stored.
2. Ask if DHCP forwarding should be enabled.
2. Ask for IP address to which DHCP captured traffic will be send to.
3. Ask for UDP  port to which DHCP captured traffic will be sent to.
4. Ask if DNS forwarding should be enabled.
5. Ask for IP address to which DNS captured traffic will be send to.
6. Ask for UDP  port to which DNS captured traffic will be sent to.
8. Store those values in Packetfence-Forwarder.toml in the working directory.
7. Store default filters value, which selects DHCPACK and DHCPREQUESTS in a DHCP mask and DNS.


Note: Do not select a wireless device, it will not work.

The Packetfence-Forwarder service needs to be restarted:

1. After a configuration change.
2. When the server goes to sleep and resumes.


The installer
-------------
The installer will:

1. Install all packaged files under "C:\Program Files (x86)\Packetfence-Forwarder"
2. Run packetfence-forwarder-config-generator.exe which generates a configuration file in installation directory.
3. Install packetfence-forwarder.exe as a service with the help of nssm.
4. Start packetfence-forwarder.exe with the help of nssm.

Build it yourself!
==================

Native Compilation under x64
----------------------------
You will need:

* [Npcap](https://npcap.com/#download)
* [Git](https://git-scm.com/download/win)
* [Go](https://golang.org/dl/)

To generate the installer, you will also need [NSIS](https://sourceforge.net/projects/nsis/files/)

Get the sources
---------------
In a shell, create a "src" directory (that is a GOLANG requirement) that you want to be the root of your GO projects, (eg: c:\Users\Test\go\src\ ) and download the sources through previously installed git:
```
git clone https://github.com/inverse-inc/packetfence-dhcp-forwarder.git
```

packetfence-forwarder-config-generator:

* Generates Packetfence-Forwarder.toml configuration based on user selected NIC from "getmac /fo csv /v"  output and fix the UID name. It is currently impossible to use gopacket to list human readable interface names so the user can choose from them and map it to its UUID.

packetfence-forwarder:

* Applies Packetfence-Forwarder.toml configuration from the working directory sends captured UDP packets to configured destination host and port.

packetfence-forwarder-installer:

* The NSI script to generate the installer is "Packetfence-Forwarder.nsi".

Files are installed under "C:\Program Files (x86)\DHCP Forwarder".

Compilation
-----------
Once you have the sources and the tools for native compilations under c:\go\src\

In a terminal, do the following:
```
cd c:\Users\Test\go\src\packetfence-dhcp-forwarder/installer
build.bat
```

You now have the compiled binaries required to generate the installer.

Compile and create the installer
--------------------
To create the installer, you need to download and install the following:
[NSIS](http://prdownloads.sourceforge.net/nsis/nsis-3.0-setup.exe?download)

Download dependency
 * Extract nssm.exe from [here](https://nssm.cc/release/nssm-2.24.zip) to c:\Users\Test\go\src\packetfence-dhcp-forwarder\installer

Place yourself in the root of the git downloaded directory:
```
cd c:\Users\Test\go\src\packetfence-packetfence-forwarder\installer
```
Compile packetfence-forwarder.exe and packetfence-forwarder-config-generator.exe and create the installer.
```
build.bat
```

You now have an installer under "c:\go\src\packetfence-packetfence-forwarder\installer\Packetfence-Forwarder-Installer.exe"


Troubleshoot
============
We officially support only x64 Windows servers.

Eventlogs
--------
The Event logs should help a lot in finding the cause of why the service not starting. Have you changed your networking card since installation? Disconnected a cable disconnected? Had the server sleep and resumed from suspend?

Alternatively, you can stop the service from Windows Service Manager and debug from the command line. Launch an adminstrative command line and place yourself under "C:\Program Files (x86)\DHCP Forwarder"

To get access to the service manager:
```
services.msc
```

To get access to event logs:
```
eventvwr.msc
```


NSSM
----
* The nssm service installation might fail if the configured interface is not in a connected state
* The nssm binary can be launched from the command line from the program files directory
* nssm configured service name is Packetfence-Forwarder

The following commands should help:

* nssm status Packetfence-Forwarder (should show a running state)
* nssm edit Packetfence-Forwarder

The service is executed with default System account. Edit accordingly.

packetfence-forwarder
--------------

* If nssm shows a status different then running, you can launch manually packetfence-forwarder from its working directory from the command line.
That application should give you more details about the reasons the service is failing.

Note: The configured interface needs to be connected. If you need to change the interface or destination information, you should execute packetfence-forwarder-config-generator.exe from a command line in its program files folder to regenerate a clean configuration.

History:
===========
* DHCP Forwarder is based on go-listener (https://github.com/louismunro/go-listener) which itself is based on the UDP reflector concept.

Installer:
* 1.0: Initial release.
* 1.1: Default "Filter" UDP port changed from 68 to 67, to make sure relays are also catched.
* 1.2: Configurator changed to not depend on english literals at installation.
* 1.3: New installator containing unified Windows and Linux shared codebase for eventlogging(google/logger).
* 1.4: Adding support for DNS forwarding rename to Packetfence-Forwarder.
