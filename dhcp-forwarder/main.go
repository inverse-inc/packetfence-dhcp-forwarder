package main

import (
	"io/ioutil"
	"os"

	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/logger"
)

func main() {
	//Setup logging
	const name = "DHCP-Forwarder"

	// Logger setup
	logger.Init(name, false, true, ioutil.Discard)
	c, err := GetConfigFromFile(name)
	checkError(err)
	handle, forwarders, err := c.SetupPcapForwarding()
	checkError(err)
	logger.Info(1, os.Args[0]+" started")
	logger.Info(1, "BPF set to: "+c.Filter)
	logger.Info(1, "Listening on device: "+c.Interface)
	//logger.Info(1, "Forwarding packets to: "+host+" on udp port "+port)
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		captureInfo := packet.Metadata().CaptureInfo
		rawPacket := packet.Data()
		for _, forwarder := range forwarders {
			if forwarder.BPF.Matches(captureInfo, rawPacket) {
				forwarder.Handler.Forward(packet)
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		logger.Error(3, err.Error())
		panic(err.Error())
	}
}
