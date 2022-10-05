package main

import (
	"net"

	"github.com/google/gopacket"
)

type ForwardHandler interface {
	Forward(p gopacket.Packet) error
}

type Forwarder struct {
	Handler ForwardHandler
	BPF     string
}

type SimpleUDPHandler struct {
	conn *net.UDPConn
}

func (h *SimpleUDPHandler) Forward(p gopacket.Packet) error {
	udpLayer := p.TransportLayer()
	if udpLayer != nil {
		conn.Write(udpLayer.LayerPayload())
		// We don't check for error here.
		// The endpoint might not be listening yet.
	}

	if errLayer := p.ErrorLayer(); errLayer != nil {
		return errLayer.Error()
	}

	return nil
}

type Config struct {
}

func NewSimpleUDPHandler(c Config) (ForwardHandler, error) {
	return &SimpleUDPHandler{}, nil
}

type DNSHandler struct {
	conn *net.UDPConn
}

func NewDNSHandler(c Config) (ForwardHandler, error) {
	return &DNSHandler{}, nil
}

func (h *DNSHandler) Forward(p gopacket.Packet) error {
	return nil
}
