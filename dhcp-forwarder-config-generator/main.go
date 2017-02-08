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

type Configuration struct {
	Host   string
	Port   int
	Device string
	Filter string
}

var Config Configuration

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
		cmdOut []byte
		err    error
	)
	//Unfortunately, gopacket device names were the same as their descriptions, so it was not possible to obtain
	//NIC's UUID directly and save it. getmac is available since WinXP, and provides the UID.
	cmdName := "getmac"
	cmdArgs := []string{"/fo", "csv", "/v"}

	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running getmac /fo csv /v", err)
		os.Exit(1)
	}
	Output := string(cmdOut)

	reader := csv.NewReader(strings.NewReader(Output))

	reader.FieldsPerRecord = 4 // Expected format is like the following: "Connection Name","Network Adapter","Physical Address","Transport Name"
	
	//We discard the header line. We expect it to not change and won't compare against all langages.
	_, err = reader.Read()
	//if CSVHeader[0] == "Connection Name" &&
	//	CSVHeader[1] == "Network Adapter" &&
	//	CSVHeader[2] == "Physical Address" &&
	//	CSVHeader[3] == "Transport Name" {
	
	rawCSVdata, err := reader.ReadAll()	
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var InterfaceIndex int
	fmt.Printf("Index\t:\t Interface name\t\n")
	for row, each := range rawCSVdata {
		fmt.Printf("%d\t:\t %s\n", row, each[1])
	}

	for {
		fmt.Printf("\nPlease select the index number corresponding to the desired interface name:")
		if _, err := fmt.Scan(&InterfaceIndex); err != nil {
			fmt.Printf("Error. %v\n", err)
		} else if 0 <= InterfaceIndex && InterfaceIndex < len(rawCSVdata) {
			//NIC's UID returned needs to be fixated by replacing Tcpip in it's name by NPF
			//NPF is WinPCAP device's driver name equivalent to the system's device.
			Config.Device = strings.Replace(rawCSVdata[InterfaceIndex][3], "Tcpip", "NPF", 1)
			break
		} else {
			fmt.Printf("!!! Choice out of possible range. Choose between 0 and %v !!!\n", len(rawCSVdata)-1)
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
