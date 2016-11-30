DHCP-Forwarder
==============

This tool captures and forwards a subset of DHCP traffic (specifically DHCPREQUEST and DHCPACK) from a Windows DHCP server to a destination IP and port.

Alternatively, IP Helpers can be configured on each switch of an infrastructure to forward broadcast only packets. Those contain all types of DHCP packets but less to none DHCPACK which confirms lease ownership.

DHCP traffic is useful to PacketFence to link MAC adresses and IP addresses, while also helping Fingerbank's fingerprinting process. Again, the only useful packets to PacketFence are DHCPREQUEST and DHCPACK.

In short, if DHCP-Forwarder can be deployed centrally, it should be done. In that case, only useful packets are captured and forwarded from the source, which reduces configuration, transport, processing and storage costs of useless packets, while supressing the need to process them to actually neglect them at the destination host.

[Download the installer here.](https://inverse.ca/downloads/PacketFence/windows-dhcp-forwarder/DHCP%20Forwarder%20Installer.exe)


Binaries
========
dhcp-forwarder.exe
------------------
This tool captures and forwards DHCP traffic (DHCPREQUEST and DHCPACK, specifically) from a Windows DHCP server to PacketFence. 

DHCPREQUEST and DHCPACK packets are the ones being the most important for PacketFence to link MAC to IP addresses and switch location and help Fingerbank to fingerprint the operating system running on the device. 

This fingerprinting and localisation process helps a lot in determining violation triggers condition.
 
With the help of DHCP-Forwarder, those DHCP packets can be obtained directly and easily from the source. 

Alternatively, IP helpers can be configured on each switch to forward DHCP traffic to PacketFence, but only broadcast packets can be captured by them, which is less precise. Deploying DHCP-Forwarder is simple and centralized.

DHCP-forwarder is based on gopacket and depends upon WinPCAP to select the requested packets through a BPF, which is really fast. Captured traffic is then forwarded to a configured host and port. 

A default BPF is produced by the configuration generator. That filter can be manually modified in the configuration file by the user.

DHCP-Forwarder requires a DHCP-Forwarder.toml file to be present in it's working directory (installation directory). DHCP-Forwarder.toml is generated from dhcp-forwarder-config-generator.exe at installation time, but can be run from the installation directory anytime.

dhcp-forwarder-config-generator.exe
-----------------------------------
Does:

1. Ask for Network Interface Card name and converts it to UUID that will be stored
2. Ask for IP address to which captured traffic will be send to
3. Ask for UDP  port to which captured traffic will be sent to.
4. Store default filter value, which selects DHCPACK and DHCPREQUESTS in a DHCP mask.
5. Store those values in DHCP-Forwarder.toml in the working directory.

Note: Do not select a Wireless device, it will not work.

The DHCP-Forwarder service needs to be restarted:

1. After a configuration change
2. When the server goes to sleep and resumes


The installer
-------------
The installer will:

1. install WinPCAP
2. install all packaged files under "C:\Program Files (x86)\DHCP Forwarder"
3. run dhcp-forwarder-config-generator.exe which generates a configuration file in installation directory
4. install dhcp-forwarder.exe as a service with the help of nssm
5. start dhcp-forwarder.exe with the help of nssm.



Build it yourself!
==================

Native Compilation under x64
----------------------------
You will need:

* [TDM-GCC](https://sourceforge.net/projects/tdm-gcc/files/latest/download)
* [WinPcap Development edition](https://sourceforge.net/projects/tdm-gcc/files/latest/download)
* [Git](https://git-scm.com/download/win)
* [Go](https://golang.org/dl/)

You will need to generate WinPCAP x64 dependencies yourself ([as of November 2016](https://stackoverflow.com/questions/38047858/compile-gopacket-on-windows-64bit). Please advise if it's not the case anymore).


To generate the installer, you will also need [NSIS](https://sourceforge.net/projects/nsis/files/)


Get the sources
---------------
In a shell, under c:\go\src or wherever you installed GO, download the sources through installed git:
```
git clone https://github.com/inverse-inc/packetfence-dhcp-forwarder.git
```
dhcp-forwarder-config-generator:

* Generates DHCP-Forwarder.toml configuration based on user selected NIC from "getmac /fo csv /v"  output and fix the UID name. It is currently impossible to use gopacket to list human readable interface names so the user can choose from them and map it to its UUID.


dhcp-forwarder:

* Applies DHCP-Forwarder.toml configuration from the working directory sends captured UDP packets to configured destination host and port.


dhcp-forwarder-installer:

* The NSI script to generate the installer is "DHCP Forwarder.nsi"

Files are installed under "C:\Program Files (x86)\DHCP Forwarder".


Compilation
-----------
Once you have the sources and the tools for native compilations under c:\go\src\

In a terminal, do the following:
```
set GOPATH=c:\go
cd c:\go\src\packetfence-dhcp-forwarder
cd dhcp-forwarder
go get
go build
cd ..
cd dhcp-forwarder-config-generator
go get
go build
```

You now have the compiled binaries required to generate the installer.


Create the installer
--------------------
To create the installer, you need to download and install the following:
[NSIS](http://prdownloads.sourceforge.net/nsis/nsis-3.0-setup.exe?download)

Place yourself in the root of the git downloaded directory:
```
cd c:\go\src\packetfence-dhcp-forwarder
```
Copy compiled files to the installer diretory:
```
cp dhcp-forwarder/dhcp-forwarder.exe dhcp-forwarder-installer
cp dhcp-forwarder-config-generator/dhcp-forwarder-config-generator.exe dhcp-forwarder-installer
```

Place yourself in the installer directory:
```
cd dhcp-forwarder-installer
```
 * extract nssm.exe from [here](https://nssm.cc/release/nssm-2.24.zip)
 * move WinPcap_4_1_3.exe to the current working directory

The following files should be present under current working directory:
 * dhcp-forwarder-installer/dhcp-forwarder-config-generator.exe
 * dhcp-forwarder-installer/DHCP Forwarder.nsi
 * dhcp-forwarder-installer/dhcp-forwarder.exe
 * dhcp-forwarder-installer/WinPcap_4_1_3.exe
 * dhcp-forwarder-installer/nssm.exe


You can now invoke the installer creator through "C:\Program Files (x86)\NSIS\NSIS.exe"
 
 * Click "Compile NSI scripts"
 * Select compression level.
 * Select "c:\go\src\packetfence-dhcp-forwarder\dhcp-forwarder-installer\DHCP Forwarder.nsi" and compile.


You now have an installer under "c:\go\src\packetfence-dhcp-forwarder\dhcp-forwarder-installer\DHCP Forwarder installer.exe"


Troubleshoot
============
Eventlogs
--------
The Event logs should help a lot in finding the cause the service not starting. Have you changed your networking card since installation? Disconnected a cable disconnected? Had the server sleep and resumed from suspend?

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
* The nssm service installation might fail if the configured interface is not in a connected state. 
* The nssm binary can be launched from the command line from the program files directory.
* nssm configured service name is DHCP-Forwarder.

The following commands should help:

* nssm status DHCP-Forwarder (should show a running state)
* nssm edit DHCP-Forwarder

The service is executed with default System account. Edit accordingly.

dhcp-forwarder
--------------

* If nssm shows a status different then running, you can launch manually dhcp-forwarder from its working directory from the command line. 
That application should give you more details about the reasons the service is failing. 

Note: The configured interface needs to be connected. If you need to change the interface or destination information, you should execute dhcp-forwarder-config-generator.exe from a command line in its program files folder to regenerate a clean configuration.


History:
===========
* DHCP Forwarder is based on go-listener (https://github.com/louismunro/go-listener) which itself is based on the UDP reflector concept.

