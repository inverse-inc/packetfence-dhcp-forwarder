package main

import (
	"io/ioutil"
	"os"
	"sync"

	_ "github.com/google/gopacket/layers"
	"github.com/google/logger"
	"github.com/inverse-inc/packetfence-dhcp-forwarder/forwarder"
)

func main() {
	//Setup logging
	const name = "PacketFence-Forwarder"

	// Logger setup
	logger.Init(name, false, true, ioutil.Discard)
	c, err := forwarder.GetConfigFromFile(name)
	checkError(err)
	forwarders, err := c.SetupPcapForwarding()
	checkError(err)
	logger.Info(1, os.Args[0]+" started")
	logger.Info(1, "BPF set to: "+c.Filter)
	logger.Info(1, "Listening on device: "+c.Interface)
	wg := &sync.WaitGroup{}
	for _, f := range forwarders {
		wg.Add(1)
		go func(f *forwarder.InterfaceForwarder) {
			defer wg.Done()
			f.HandlePackets()
			//logger.Info(1, "Forwarding packets to: "+host+" on udp port "+port)
		}(f)
	}

	wg.Wait()
}

func checkError(err error) {
	if err != nil {
		logger.Error(3, err.Error())
		panic(err.Error())
	}
}
