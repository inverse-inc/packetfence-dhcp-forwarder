package main

import (
	"io/ioutil"
	"net"
	"os"

	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strconv"

	"syscall"

	"github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/routing"
	"github.com/google/logger"
	"github.com/mdlayher/raw"
	"github.com/spf13/viper"
)

var (
	confDirs = [3]string{"/etc/DHCP-Forwarder/", "$HOME/.DHCP-Forwarder", "."}

	filter  string
	exclude string
	conn    *net.UDPConn
	dev     string
	snaplen int32 = 1600
	host    string
	port    string
	srcIP   *net.IP
)

const UDP_HEADER_LEN = 8

// A RawClient is a Wake-on-LAN client which operates directly on top of
// Ethernet frames using raw sockets.  It can be used to send WoL magic packets
// to other machines on a local network, using their hardware addresses.
type RawClient struct {
	ifi *net.Interface
	p   net.PacketConn
}

type udphdr struct {
	src  uint16
	dst  uint16
	ulen uint16
	csum uint16
}

type pseudohdr struct {
	ipsrc   [4]byte
	ipdst   [4]byte
	zero    uint8
	ipproto uint8
	plen    uint16
}

type iphdr struct {
	vhl   uint8
	tos   uint8
	iplen uint16
	id    uint16
	off   uint16
	ttl   uint8
	proto uint8
	csum  uint16
	src   [4]byte
	dst   [4]byte
}

func main() {
	//Setup logging
	const name = "DHCP-Forwarder"

	// Logger setup
	logger.Init(name, false, true, ioutil.Discard)

	viper.SetConfigName("DHCP-Forwarder") // will match DHCP-Forwarder.{toml,json} etc.

	pwd, err := os.Getwd()
	checkError(err)

	viper.AddConfigPath(pwd)

	err = viper.ReadInConfig()
	checkError(err)

	host = viper.GetString("DestinationHost")
	port = viper.GetString("DestinationPort")
	dev = viper.GetString("ListeningDevice")
	filter = viper.GetString("Filter")
	exclude = " and (not (dst port " + port + " and dst host " + host + " ))"

	udpAddr, err := net.ResolveUDPAddr("udp4", host+":"+port)
	checkError(err)

	conn, err = net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	handle, err := pcap.OpenLive(dev, snaplen, true, pcap.BlockForever)
	checkError(err)

	route, err := routing.New()
	checkError(err)

	_, _, srcIP, err := route.Route(net.ParseIP(host))
	checkError(err)

	err = handle.SetBPFFilter(filter + exclude)
	checkError(err)

	logger.Info(1, os.Args[0]+" started")
	logger.Info(1, "BPF set to: "+filter+exclude)
	logger.Info(1, "Listening on device: "+dev)
	logger.Info(1, "Forwarding packets to: "+host+" on udp port "+port)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		handlePacket(packet, srcIP)
	}
}

func handlePacket(p gopacket.Packet, srcIP net.IP) {
	udpLayer := p.TransportLayer()
	if udpLayer != nil {
		portint, _ := strconv.Atoi(port)

		sendUnicast(udpLayer.LayerPayload(), net.UDPAddr{net.ParseIP(host), portint, ""}, srcIP)

		// We don't check for error here.
		// The endpoint might not be listening yet.
	}
	if err := p.ErrorLayer(); err != nil {
		logger.Info(2, "Error decoding some part of the packet.")
	}
}

func checkError(err error) {
	if err != nil {
		logger.Error(3, err.Error())
		panic(err.Error())
	}
}

func sendUnicast(udpPayload []byte, dstIP net.UDPAddr, srcIP net.IP) error {

	s, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Fatal(err)
	}

	proto := 17

	ipStr, portStr, _ := net.SplitHostPort(dstIP.String())
	port, _ := strconv.Atoi(portStr)

	udpsrc := uint(67)
	udpdst := port

	udp := udphdr{
		src: uint16(udpsrc),
		dst: uint16(udpdst),
	}

	udplen := 8 + len(udpPayload)

	ip := iphdr{
		vhl:   0x45,
		tos:   0,
		id:    0x0000, // the kernel overwrites id if it is zero
		off:   0,
		ttl:   128,
		proto: uint8(proto),
	}
	copy(ip.src[:], net.ParseIP(srcIP.String()).To4())
	copy(ip.dst[:], net.ParseIP(ipStr).To4())

	udp.ulen = uint16(udplen)
	udp.checksum(&ip, udpPayload)

	totalLen := 20 + udplen

	ip.iplen = uint16(totalLen)
	ip.checksum()

	buf := bytes.NewBuffer([]byte{})
	err = binary.Write(buf, binary.BigEndian, &udp)
	if err != nil {
		log.Fatal(err)
	}

	udpHeader := buf.Bytes()
	dataWithHeader := append(udpHeader, udpPayload...)

	buff := bytes.NewBuffer([]byte{})
	err = binary.Write(buff, binary.BigEndian, &ip)
	if err != nil {
		log.Fatal(err)
	}

	ipHeader := buff.Bytes()
	packet := append(ipHeader, dataWithHeader...)

	addr := syscall.SockaddrInet4{}
	copy(addr.Addr[:], net.ParseIP(ipStr).To4())
	addr.Port = int(udpdst)

	err = syscall.Sendto(s, packet, 0, &addr)
	// Send packet to target
	err = syscall.Close(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error closing the socket: %v\n", err)
		os.Exit(1)
	}
	return err
}

// NewRawClient creates a new RawClient using the specified network interface.
//
// Note that raw sockets typically require elevated user privileges, such as
// the 'root' user on Linux, or the 'SET_CAP_RAW' capability.
//
// For this reason, it is typically recommended to use the regular Client type
// instead, which operates over UDP.
func NewRawClient(ifi *net.Interface) (*RawClient, error) {
	// Open raw socket to send Wake-on-LAN magic packets
	var cfg raw.Config

	p, err := raw.ListenPacket(ifi, 0x0806, &cfg)
	if err != nil {
		return nil, err
	}

	return &RawClient{
		ifi: ifi,
		p:   p,
	}, nil
}

// Close closes a RawClient's raw socket.
func (c *RawClient) Close() error {
	return c.p.Close()
}

func (u *udphdr) checksum(ip *iphdr, payload []byte) {
	u.csum = 0

	phdr := pseudohdr{
		ipsrc:   ip.src,
		ipdst:   ip.dst,
		zero:    0,
		ipproto: ip.proto,
		plen:    u.ulen,
	}
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, &phdr)
	binary.Write(&b, binary.BigEndian, u)
	binary.Write(&b, binary.BigEndian, &payload)
	u.csum = checksum(b.Bytes())
}

func checksum(buf []byte) uint16 {
	sum := uint32(0)

	for ; len(buf) >= 2; buf = buf[2:] {
		sum += uint32(buf[0])<<8 | uint32(buf[1])
	}
	if len(buf) > 0 {
		sum += uint32(buf[0]) << 8
	}
	for sum > 0xffff {
		sum = (sum >> 16) + (sum & 0xffff)
	}
	csum := ^uint16(sum)
	/*
	 * From RFC 768:
	 * If the computed checksum is zero, it is transmitted as all ones (the
	 * equivalent in one's complement arithmetic). An all zero transmitted
	 * checksum value means that the transmitter generated no checksum (for
	 * debugging or for higher level protocols that don't care).
	 */
	if csum == 0 {
		csum = 0xffff
	}
	return csum
}

func (h *iphdr) checksum() {
	h.csum = 0
	var b bytes.Buffer
	binary.Write(&b, binary.BigEndian, h)
	h.csum = checksum(b.Bytes())
}
