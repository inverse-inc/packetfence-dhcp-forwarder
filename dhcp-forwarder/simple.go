package main

import (
	"net"

	"github.com/google/gopacket"
)

type SimpleUDPHandler struct {
	conn *net.UDPConn
}

func (h *SimpleUDPHandler) Forward(p gopacket.Packet) error {
	udpLayer := p.TransportLayer()
	if udpLayer != nil {
		h.conn.Write(udpLayer.LayerPayload())
		// We don't check for error here.
		// The endpoint might not be listening yet.
	}

	if errLayer := p.ErrorLayer(); errLayer != nil {
		return errLayer.Error()
	}

	return nil
}

func NewSimpleUDPHandler(c *ForwarderConfig) (ForwardHandler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &SimpleUDPHandler{conn: conn}, nil
}
