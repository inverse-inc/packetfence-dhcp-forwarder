package main

import (
	"errors"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type ForwardHandler interface {
	Forward(p gopacket.Packet) error
}

type Forwarder struct {
	Handler ForwardHandler
	BPF     *pcap.BPF
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

func NewSimpleUDPHandler(c *ForwarderConfig) (ForwardHandler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &SimpleUDPHandler{conn: conn}, nil
}

type DNSHandler struct {
	conn *net.UDPConn
}

func NewDNSHandler(c *ForwarderConfig) (ForwardHandler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", c.Host+":"+c.Port)
	if err != nil {
		return nil, err
	}

	conn, err = net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &DNSHandler{conn: conn}, nil
}

func (h *DNSHandler) Forward(p gopacket.Packet) error {
	return nil
}

type ForwardHandlerBuilder func(c *ForwarderConfig) (ForwardHandler, error)

var forwardHandlerBuilders = map[string]ForwardHandlerBuilder{
	"dns":  NewDNSHandler,
	"dhcp": NewSimpleUDPHandler,
}

func MakeForwarder(pcap *pcap.Handle, c *ForwarderConfig) (*Forwarder, error) {
	if f, found := forwardHandlerBuilders[c.Type]; found {
		h, err := f(c)
		if err != nil {
			return nil, err
		}

		bpf, err := pcap.NewBPF(c.Filter)
		if err != nil {
			return nil, err
		}

		return &Forwarder{Handler: h, BPF: bpf}, nil
	}

	return nil, errors.New("Not found")
}
