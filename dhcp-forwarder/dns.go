package main

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type DNSHandler struct {
	conn *net.UDPConn
}

func NewDNSHandler(c *ForwarderConfig) (ForwardHandler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &DNSHandler{conn: conn}, nil
}

func (h *DNSHandler) Forward(p gopacket.Packet) error {
	var eth layers.Ethernet
	var ipv4 layers.IPv4
	var udp layers.UDP
	var dns layers.DNS
	parser := gopacket.NewDecodingLayerParser(layers.LayerTypeEthernet, &eth, &ipv4, &udp, &dns)
	decoded := make([]gopacket.LayerType, 0, 4)
	if err := parser.DecodeLayers(p.Data(), &decoded); err != nil {
		return err
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := dns.SerializeTo(buf, opts); err != nil {
		return err
	}

	h.conn.Write(buf.Bytes())
	return nil
}
