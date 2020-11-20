package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/gopacket/pcap"
	"golang.org/x/sys/windows/registry"
)

type Configuration struct {
	DestinationHost string
	DestinationPort int
	ListeningDevice string
	Filter          string
}

var Config Configuration

type NetworkInterface struct {
	Name        string
	Description string
}

func main() {
	fmt.Printf("!!! You can run this program anytime from %s !!!\n\n", os.Args[0])
	//Set values to defaults. 0x3: DHCPREQUEST. 0x5: DHCPACK. Those are the only ones required by PacketFence to track and fingerprint devices from DHCP.
	Config.Filter = "udp and port 67 and ((udp[250:1] = 0x3) or (udp[250:1] = 0x5))"
	SelectInterface()
	SelectRemoteHostAndPort()
	SaveConfig("DHCP-Forwarder.toml")
}

func SelectInterface() {
	var (
		err error
	)

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	interfacePattern := regexp.MustCompile("\\{(.*)\\}")
	NetworkInterfaces := []NetworkInterface{}

	for _, device := range devices {
		NetInterface := &NetworkInterface{}
		NetInterface.Name = device.Name
		match := interfacePattern.FindStringSubmatch(strings.ToLower(device.Name))
		k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\ControlSet001\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}\`+match[0]+`\Connection`, registry.QUERY_VALUE)
		if err != nil {
			log.Println(err)
		}
		defer k.Close()
		s, _, err := k.GetStringValue("Name")
		if err != nil {
			continue
		}
		NetInterface.Description = s
		NetworkInterfaces = append(NetworkInterfaces, *NetInterface)
	}

	var InterfaceIndex int
	fmt.Printf("Index\t:\t Interface name\t\n")
	for row, each := range NetworkInterfaces {
		fmt.Printf("%d\t:\t %s\n", row, each.Description)
	}

	for {
		fmt.Printf("\nPlease select the index number corresponding to the desired interface name:")
		if _, err := fmt.Scan(&InterfaceIndex); err != nil {
			fmt.Printf("Error. %v\n", err)
		} else if 0 <= InterfaceIndex && InterfaceIndex < len(NetworkInterfaces) {
			//NIC's UID returned needs to be fixated by replacing Tcpip in it's name by NPF
			//NPF is WinPCAP device's driver name equivalent to the system's device.
			Config.ListeningDevice = NetworkInterfaces[InterfaceIndex].Name
			break
		} else {
			fmt.Printf("!!! Choice out of possible range. Choose between 0 and %v !!!\n", len(NetworkInterfaces)-1)
		}
	}
}

func SelectRemoteHostAndPort() {
	var tmp string
	var UDPPort int

	for {
		fmt.Printf("To which IP will the selected UDP traffic be forwarded to? ")
		if _, err := fmt.Scan(&tmp); err != nil {
			fmt.Printf("Error. %v\n\n", err)
		}
		TestInput := net.ParseIP(tmp)
		if TestInput.To4() == nil {
			fmt.Printf("!!! %v is not a valid hostv4 address !!!\n", tmp)
		} else {
			Config.DestinationHost = tmp
			break
		}
	}

	for {
		fmt.Printf("To which UDP port will the selected UDP traffic be forwarded to? ")
		if _, err := fmt.Scan(&UDPPort); err != nil {
			fmt.Printf("Error. %v\n\n", err)
		} else if 0 <= UDPPort && UDPPort <= 65535 {
			Config.DestinationPort = UDPPort
			break
		} else {
			fmt.Printf("!!! UDP port out of possible range. Choose between 0 and 65535 !!!\n")
		}
	}
}

func SaveConfig(configFile string) {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	encoder := toml.NewEncoder(writer)
	err := encoder.Encode(Config)

	if err != nil {
		fmt.Printf("Error. %s\n", err)
	}

	if err := writer.Flush(); err != nil {
		fmt.Printf("Error. %s\n", err)
	}

	if err := ioutil.WriteFile(configFile, buf.Bytes(), 0644); err != nil {
		fmt.Printf("Error. %s\n", err)
	}
}
