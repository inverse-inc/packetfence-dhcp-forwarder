package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	Host   string
	Port   int
	Device string
	Filter string
}

var Config TomlConfig

func main() {
	fmt.Printf("!!! You can run this program anytime from %s !!!\n\n", os.Args[0])
	
	//Set values to defaults. 0x3: DHCPREQUEST. 0x5: DHCPACK. Those are the only ones required by PacketFence to track and fingerprint devices from DHCP.
	Config.Filter = "udp and port 68 and ((udp[250:1] = 0x3) or (udp[250:1] = 0x5))"
	SelectInterface()
	SelectRemoteHostAndPort()
	
	//fmt.Printf("Actual configuration:\n")
	//fmt.Printf("Config.filter\t%v\n", Config.Filter)
	//fmt.Printf("Config.host\t%v\n", Config.Host)
	//fmt.Printf("Config.port\t%v\n", Config.Port)
	//fmt.Printf("Config.device\t%v\n", Config.Device)
		
	SaveConfig("DHCP-Forwarder.toml")
}

func SelectInterface() {
	var (
		cmdOut []byte
		err    error
	)

	cmdName := "getmac"
	cmdArgs := []string{"/fo", "csv", "/v"}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running getmac /fo csv /v", err)
		os.Exit(1)
	}
	Output := string(cmdOut)

	reader := csv.NewReader(strings.NewReader(Output))

	reader.FieldsPerRecord = 4 // Expected format is "Connection Name","Network Adapter","Physical Address","Transport Name"
	CSVHeader, err := reader.Read()

	//If the ouput of getmac is in the expected format
	if CSVHeader[0] == "Connection Name" &&
		CSVHeader[1] == "Network Adapter" &&
		CSVHeader[2] == "Physical Address" &&
		CSVHeader[3] == "Transport Name" {
		rawCSVdata, err := reader.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		var InterfaceIndex int
		fmt.Printf("Index\t:\t Interface name\t\n")
		for row, each := range rawCSVdata {
			fmt.Printf("%d\t:\t %s\n", row, each[0])
		}

		for {
			fmt.Printf("\nPlease select the index number corresponding to the desired interface name:")
			if _, err := fmt.Scan(&InterfaceIndex); err != nil {
				fmt.Printf("Error. %v\n", err)
			} else if 0 <= InterfaceIndex && InterfaceIndex < len(rawCSVdata) {
				Config.Device = strings.Replace(rawCSVdata[InterfaceIndex][3], "Tcpip", "NPF", 1)
				break
			} else {
				fmt.Printf("!!! Choice out of possible range. Choose between 0 and %v !!!\n", len(rawCSVdata)-1)
			}
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
			Config.Host = tmp
			break
		}
	}

	for {
		fmt.Printf("To which UDP port will the selected UDP traffic be forwarded to? ")
		if _, err := fmt.Scan(&UDPPort); err != nil {
			fmt.Printf("Error. %v\n\n", err)
		} else if 0 <= UDPPort && UDPPort <= 65535 {
			Config.Port = UDPPort
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