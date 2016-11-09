package main

import (
	"fmt"
	"net"
	"os"
	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/viper"
)

var (
	confDirs = [3]string{"/etc/DHCP-Forwarder/", "$HOME/.DHCP-Forwarder", "."}
	
	filter   string
	exclude  string
	conn     *net.UDPConn
	dev      string
	snaplen  int32 = 1600
	host     string
	port     string
)

func main() {

	//viper.SetDefault("host", "localhost")
	//viper.SetDefault("port", "6767")
	//viper.SetDefault("device", "none")
	//viper.SetDefault("filter", "udp and port 68 and ((udp[250:1] = 0x3) or (udp[250:1] = 0x5))")

	viper.SetConfigName("DHCP-Forwarder") // will match DHCP-Forwarder.{toml,json} etc.
	
	pwd, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Println("PWD: ", pwd)
	
	//for _, dir := range confDirs {
	//	viper.AddConfigPath(dir)
	//}
	viper.AddConfigPath(pwd)
	
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	host = viper.GetString("Host")
	port = viper.GetString("Port")
	dev = viper.GetString("Device")
	filter = viper.GetString("Filter")
	exclude = " and (not (dst port " + port + " and dst host " + host + " ))"

	udpAddr, err := net.ResolveUDPAddr("udp4", host+":"+port)
	checkError(err)

	conn, err = net.DialUDP("udp", nil, udpAddr)
	checkError(err)	
	
	if handle, err := pcap.OpenLive(dev, snaplen, true, pcap.BlockForever); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(filter + exclude); err != nil {
		panic(err)
	} else {
		fmt.Println(os.Args[0] + " started")
		fmt.Println("Listening on device " + dev)
		fmt.Println("BPF set to " + filter + exclude)
		fmt.Println("Forwarding packets to " + host + ":" + port)

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			handlePacket(packet)
		}
	}
}

func handlePacket(p gopacket.Packet) {
	udpLayer := p.TransportLayer()
	if udpLayer != nil {
		conn.Write(udpLayer.LayerPayload())
		// We don't check for error here.
		// The endpoint might not be listening yet.
	}
	if err := p.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}
