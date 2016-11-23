package main

import (
	"net"
	"os"
	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/spf13/viper"
	"golang.org/x/sys/windows/svc/eventlog"
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
	
	Logger *eventlog.Log
)

func main() {
    //Setup logging
	const name = "DHCP-Forwarder"
	const supports = eventlog.Error | eventlog.Warning | eventlog.Info
	eventlog.Remove(name)
	err := eventlog.InstallAsEventCreate(name, supports)
	checkError(err)
	defer func() {
		err = eventlog.Remove(name)
		checkError(err)
	}()
	
	Logger, err = eventlog.Open(name)
	checkError(err)
	defer Logger.Close()
	
	viper.SetConfigName("DHCP-Forwarder") // will match DHCP-Forwarder.{toml,json} etc.
	
	pwd, err := os.Getwd()
    checkError(err)

	viper.AddConfigPath(pwd)
	
	err = viper.ReadInConfig()
	checkError(err)
	
	host = viper.GetString("Host")
	port = viper.GetString("Port")
	dev = viper.GetString("Device")
	filter = viper.GetString("Filter")
	exclude = " and (not (dst port " + port + " and dst host " + host + " ))"

	udpAddr, err := net.ResolveUDPAddr("udp4", host+":"+port)
	checkError(err)

	conn, err = net.DialUDP("udp", nil, udpAddr)
	checkError(err)
	
	handle, err := pcap.OpenLive(dev, snaplen, true, pcap.BlockForever)
	checkError(err)
	
	err = handle.SetBPFFilter(filter + exclude)
	checkError(err)
	
	Logger.Info(1, os.Args[0] + " started")
	Logger.Info(1, "BPF set to: " + filter + exclude)
	Logger.Info(1, "Listening on device: " + dev)
	Logger.Info(1, "Forwarding packets to: " + host + " on udp port " + port)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handlePacket(packet)
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
		Logger.Warning(2, "Error decoding some part of the packet.")
	}
}

func checkError(err error) {
	if err != nil {
		Logger.Error(3, err.Error())
		panic(err.Error())
	}
}
