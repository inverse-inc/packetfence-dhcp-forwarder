package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/google/gopacket/pcap"
	"golang.org/x/sys/windows/registry"
)

func getInterfaces() []NetworkInterface {
	interfacePattern := regexp.MustCompile("\\{(.*)\\}")
	networkInterfaces := []NetworkInterface{}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for _, device := range devices {
		netInterface := NetworkInterface{}
		netInterface.Name = device.Name
		match := interfacePattern.FindStringSubmatch(strings.ToLower(device.Name))
		if len(match) > 0 {
			k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}\`+match[0]+`\Connection`, registry.QUERY_VALUE)
			if err != nil {
				log.Println(err)
			}
			s, _, err := k.GetStringValue("Name")
			k.Close()
			if err != nil {
				continue
			}

			netInterface.Description = s
		} else {
			netInterface.Description = device.Name
		}

		networkInterfaces = append(networkInterfaces, netInterface)
	}

	return networkInterfaces
}
