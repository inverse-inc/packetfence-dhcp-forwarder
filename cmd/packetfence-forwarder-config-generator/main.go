package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Configuration struct {
	ListeningDevice string

	DisableDHCP     bool
	DestinationHost string
	DestinationPort uint16
	Filter          string

	DisableDNS         bool
	DNSDestinationHost string
	DNSDestinationPort uint16
	DNSFilter          string
}

func defaultConfig() Configuration {
	return Configuration{
		DisableDHCP:        true,
		DestinationHost:    "",
		DestinationPort:    767,
		Filter:             "udp and port 67 and ((udp[250:1] = 0x3) or (udp[250:1] = 0x5))",
		DisableDNS:         true,
		DNSDestinationHost: "",
		DNSDestinationPort: 753,
		DNSFilter:          "udp and port 53",
	}
}

func main() {
	fmt.Printf("!!! You can run this program anytime from %s !!!\n\n", os.Args[0])
	scanner := bufio.NewScanner(os.Stdin)
	//Set values to defaults. 0x3: DHCPREQUEST. 0x5: DHCPACK. Those are the only ones required by PacketFence to track and fingerprint devices from DHCP.
	config := defaultConfig()
	SelectInterface(scanner, &config)
	SetupDHCPForwarding(scanner, &config)
	SetupDNSForwarding(scanner, &config)
	SaveConfig(&config, "PacketFence-Forwarder.toml")
}

func SelectInterface(scanner *bufio.Scanner, c *Configuration) {

	networkInterfaces := getInterfaces()

	var interfaceChoose string
	fmt.Printf("Index\t:\t Interface name\t\n")
	for row, each := range networkInterfaces {
		fmt.Printf("%d\t:\t %s\n", row, each.Description)
	}

	for {
		fmt.Printf("\nPlease choose the number interface or press enter for all: ")
		scanner.Scan()
		interfaceChoose = scanner.Text()

		if interfaceChoose == "all" || interfaceChoose == "" {
			c.ListeningDevice = ""
			break
		}

		interfaceIndex, err := strconv.ParseInt(interfaceChoose, 10, 64)
		if err == nil {
			if 0 <= interfaceIndex && interfaceIndex < int64(len(networkInterfaces)) {
				c.ListeningDevice = networkInterfaces[interfaceIndex].Name
				break
			}
		}

		fmt.Printf("!!! Choice out of possible range. Choose between 0 and %v or all !!!\n", len(networkInterfaces)-1)
	}
}

func setupEnable(scanner *bufio.Scanner, msg string) bool {
	for {
		fmt.Printf("\n%s y/n: ", msg)
		scanner.Scan()
		enable := scanner.Text()

		if enable == "" {
			continue
		}

		if enable[0] == 'Y' || enable[0] == 'y' {
			return true
		}

		if enable[0] == 'N' || enable[0] == 'n' {
			return false
		}

	}
}

func SetupDHCPForwarding(scanner *bufio.Scanner, c *Configuration) {
	if !setupEnable(scanner, "Enable DHCP forwarding") {
		c.DisableDHCP = true
		return
	}

	c.DisableDHCP = false
	setupHostAndPort(scanner, "DHCP", &c.DestinationHost, &c.DestinationPort, 767)
}

func SetupDNSForwarding(scanner *bufio.Scanner, c *Configuration) {
	if !setupEnable(scanner, "Enable DNS forwarding") {
		c.DisableDNS = true
		return
	}

	c.DisableDNS = false
	setupHostAndPort(scanner, "DNS", &c.DNSDestinationHost, &c.DNSDestinationPort, 753)
}

func setupHostAndPort(scanner *bufio.Scanner, msg string, host *string, port *uint16, defaultPort uint16) {
	fmt.Printf("\nSetting up %s forward host and port\n", msg)
	for {
		fmt.Printf("To which IP will the selected %s traffic be forwarded to: ", msg)
		tmpStr := ""
		for tmpStr == "" {
			scanner.Scan()
			tmpStr = scanner.Text()
		}

		testInput := net.ParseIP(tmpStr)
		if testInput.To4() == nil {
			fmt.Printf("!!! %v is not a valid hostv4 address !!!\n", tmpStr)
			continue
		} else {
			*host = tmpStr
			break
		}
	}

	*port = defaultPort
	for {
		fmt.Printf("To which UDP port will the selected %s traffic be forwarded to press enter for the default port %d: ", msg, defaultPort)
		scanner.Scan()
		tmpStr := scanner.Text()
		if tmpStr == "" {
			break
		}

		tmpPort, err := strconv.ParseInt(tmpStr, 10, 16)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			continue
		}

		if 0 < tmpPort && tmpPort <= 65535 {
			*port = uint16(tmpPort)
			break
		}

		fmt.Printf("!!! UDP port out of possible range. Choose between 1 and 65535 !!!\n")
	}
}

func SaveConfig(c *Configuration, configFile string) {
	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)

	encoder := toml.NewEncoder(writer)
	err := encoder.Encode(c)

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
