The installer
===========
The installer will:

1. Install WinPCAP
2. Install all packaged files under "C:\Program Files (x86)\DHCP Forwarder"
3. Run dhcp-forwarder-config-generator.exe which generates a configuration file in previous directory.
4. Install dhcp-forwarder.exe as a service with the help of nssm
5. Start dhcp-forwarder.exe with the help of nssm.

The installer is downloadable from: 
dhcp-forwarder-installer\DHCP Forwarder Installer.exe

dhcp-forwarder-installer\* contains all binaries packaged in the installer.


dhcp-forwarder-config-generator.exe
===========
Does:

1. Ask for Network Interface Card name and converts it to UUID that will be stored
2. Ask for IP address to which captured traffic will be send to
3. Ask for UDP  port to which captured traffic will be sent to.
4. Store default filter value, which selects DHCPACK and DHCPREQUESTS in a DHCP mask.
5. Store those values in DHCP-Forwarder.toml in the working directory.


dhcp-forwarder.exe
===========
A quick go tool to mirror a UDP DHCP packet capture based on a BPF and forward it to a configured host and port.
Support a configuration file in DHCP-Forwarder in TOML format, that needs to be present in the same directory as the binary.

DHCP-Forwarder.toml file is expected to be present in the working directory, which is generated from dhcp-forwarder-config-generator.exe




Build it yourself!
===========

Native Compilation under x64:
===========
You will need:

* TDM-GCC: https://sourceforge.net/projects/tdm-gcc/files/latest/download
* WinPcap Development edition: https://sourceforge.net/projects/tdm-gcc/files/latest/download
* Git: https://git-scm.com/download/win
* Go: https://golang.org/dl/

You will need to generate WinPCAP x64 dependencies yourself (as of November 2016. Please advise if it's not the case anymore.

https://stackoverflow.com/questions/38047858/compile-gopacket-on-windows-64bit

To generate the installer, you will also need NSIS: https://sourceforge.net/projects/nsis/files/


Get the sources:
===========
In a shell, under c:\go\src or wherever you installed GO, download the sources through installed git:

git clone https://github.com/inverse-inc/packetfence-dhcp-forwarder.git

dhcp-forwarder-config-generator:
Generates DHCP-Forwarder.toml configuration based on user selection of modified "getmac /fo csv /v" output, since it is currently impossible to use gopacket to list human readable interface name.
Asks for destination ip and port.
The configuration file contains preconfigured BPF for DHCPREQUESTS(3) and DHCPACK(5). That can be modified directly in the configuration file.

Notes: Wireless device selection will not work.


dhcp-forwarder:
Takes DHCP-Forwarder.toml from the working directory and sends captured udp packet to configured destination host and port.


dhcp-forwarder-installer:
The installer will launch WinPCAP Installer, install all required files (dhcp-forwarder-config-generator.exe, dhcp-forwarder.exe, nssm to manage services) and launch the ConfigGenerator and install dhcp-forwarder.exe as a service with nssm.
The NSI script is DHCP Forwarder.nsi

Files are installed under "C:\Program Files (x86)\DHCP Forwarder".

Compilation:
===========
Once you have the sources and the tools for Native compilations under c:\go\src\

In a terminal, do the following:
set GOPATH=c:\go

cd c:\go\src\packetfence-dhcp-forwarder

cd dhcp-forwarder
go get
go build
cd ..
cd dhcp-forwarder-config-generator
go get
go build

You now have the binaries required to build the installer.


Create the installer:
==========
To create the installer, you need to download and install the following:
NSIS: http://prdownloads.sourceforge.net/nsis/nsis-3.0-setup.exe?download

Place yourself in the root of the git downloaded directory:
cd c:\go\src\packetfence-dhcp-forwarder

Copy compiled files to the installer diretory:
cp dhcp-forwarder/dhcp-forwarder.exe dhcp-forwarder-installer
cp dhcp-forwarder-config-generator/dhcp-forwarder-config-generator.exe dhcp-forwarder-installer


Place yourself in the installer directory:
cd dhcp-forwarder-installer
extract nssm.exe from https://nssm.cc/release/nssm-2.24.zip
move WinPcap_4_1_3.exe here.

The following files should be present under current working directory:
dhcp-forwarder-installer/dhcp-forwarder-config-generator.exe
dhcp-forwarder-installer/DHCP Forwarder.nsi
dhcp-forwarder-installer/dhcp-forwarder.exe
dhcp-forwarder-installer/WinPcap_4_1_3.exe
dhcp-forwarder-installer/nssm.exe


You can invoke the installer creator through "C:\Program Files (x86)\NSIS\NSIS.exe"
Click "Compile NSI scripts", select "c:\go\src\packetfence-dhcp-forwarder\dhcp-forwarder-installer\DHCP Forwarder.nsi" and compile.
Select compression level.
You now have an installer under "c:\go\src\packetfence-dhcp-forwarder\dhcp-forwarder-installer\DHCP Forwarder installer.exe"



Troubleshoot:
===========
Read the eventlogs, this should help a lot in finding the cause of a service not starting. Have you changed your networking card connection of got a cable disconnected?


Launch an adminstrative command line and place yourself under "C:\Program Files (x86)\DHCP Forwarder"

NSSM:

* The nssm service installation might fail if the configured interface is not in a connected state. 
* The nssm binary can be launched from the command line from the program files directory.
nssm configured service name is DHCP-Forwarder.

The following commands should help:

* nssm status DHCP-Forwarder (should show a running state)
* nssm edit DHCP-Forwarder

The service is run with default System account. Edit accordingly.

dhcp-forwarder:

* If nssm shows a status different then running, you can launch manually dhcp-forwarder from its working directory from the command line. 
That application should give you more details about the reasons the service is failing. 

Note: The configured interface needs to be connected. If you need to change the interface or destination information, you should rerun dhcp-forwarder-config-generator.exe from a command line in its program files folder.


History:
===========
* DHCP Forwarder is based on go-listener (https://github.com/louismunro/go-listener) which itself is based on the udp reflector concept.

