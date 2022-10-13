package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket/pcap"
)

func getInterfaces() []NetworkInterface {
	networkInterfaces := []NetworkInterface{}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		if len(device.Addresses) == 0 {
			continue
		}

		netInterface := NetworkInterface{}
		netInterface.Name = device.Name
		if device.Description != "" {
			netInterface.Description = device.Description
		} else {
			netInterface.Description = device.Name
		}

		fmt.Printf("%08b: %s\n", device.Flags, device.Name)
		networkInterfaces = append(networkInterfaces, netInterface)
	}

	return networkInterfaces
}
