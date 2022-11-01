packetfence-forwarder-installer
========================

> Only used for Microsoft Windows based system, Linux based system instructions can be found [here](https://github.com/inverse-inc/packetfence-packetfence-forwarder/tree/master/packetfence-forwarder#linux-based-system)

Part of the PacketFence-Forwarder, 'packetfence-forwarder-installer' helps Microsoft Windows based system step-by-step installation.

 * The NSI script needed to generate the installer is "Packetfence-Forwarder.nsi"
 * Files are installed under "C:\Program Files (x86)\Packetfence Forwarder".

Creating the Installer
========================

Install [NSIS](https://nsis.sourceforge.io/Download) compiler.
Download [nsse.exe] https://nssm.cc/download and place it here.

Compiling
========================
Download golang https://go.dev/dl/

https://go.dev/dl/go1.18.7.windows-amd64.msi

cd cmd/
packetfence-forwarder  packetfence-forwarder-config-generator
