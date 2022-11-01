cd ..\cmd\packetfence-forwarder
"c:\Program Files\Go\bin\go.exe" build
copy packetfence-forwarder.exe ..\..\installer

cd ..\packetfence-forwarder-config-generator
"c:\Program Files\Go\bin\go.exe" build
copy packetfence-forwarder-config-generator.exe ..\..\installer

cd ..\..\installer
"c:\Program Files (x86)\NSIS\makensis.exe" /X"SetCompressor /FINAL lzma" Packetfence-Forwarder.nsi