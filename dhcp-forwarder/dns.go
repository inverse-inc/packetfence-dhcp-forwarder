package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type DNSHandler struct {
	conn *net.UDPConn
	addr *net.UDPAddr
}

func NewDNSHandler(c *ForwarderConfig) (ForwardHandler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
	if err != nil {
		return nil, err
	}

	return &DNSHandler{conn: conn, addr: udpAddr}, nil
}

func (h *DNSHandler) Forward(p gopacket.Packet) error {
	var eth layers.Ethernet
	var ipv4 layers.IPv4
	var ipv6 layers.IPv6
	var udp layers.UDP
	var dns layers.DNS
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ipv4, &ipv6, &udp, &dns)
	decoded := make([]gopacket.LayerType, 0, 4)
	if err := parser.DecodeLayers(p.Data(), &decoded); err != nil {
		return err
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true}
	if err := dns.SerializeTo(buf, opts); err != nil {
		return err
	}

	h.conn.WriteTo(buf.Bytes(), h.addr)
	return nil
}
