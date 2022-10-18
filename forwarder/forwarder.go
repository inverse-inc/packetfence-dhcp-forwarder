package forwarder

import (
	"errors"

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
