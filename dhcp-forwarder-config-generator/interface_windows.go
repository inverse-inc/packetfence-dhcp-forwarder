package main

import (
	"fmt"
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
		netInterface = NetworkInterface{}
		netInterface.Name = device.Name
		fmt.Printf("%s\n", device.Name)
		networkInterfaces = append(NetworkInterfaces, NetInterface)
		match := interfacePattern.FindStringSubmatch(strings.ToLower(device.Name))
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}\`+match[0]+`\Connection`, registry.QUERY_VALUE)
		if err != nil {
			log.Println(err)
		}
		s, _, err := k.GetStringValue("Name")
		if err != nil {
			continue
		}

		NetInterface.Description = s
		k.Close()
	}

	return networkInterfaces
}
