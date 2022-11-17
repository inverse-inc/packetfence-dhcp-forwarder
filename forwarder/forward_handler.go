package forwarder

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type InterfaceForwarder struct {
	handle     *pcap.Handle
	forwarders []*Forwarder
}

func (fh *InterfaceForwarder) HandlePackets() {
	packetSource := gopacket.NewPacketSource(fh.handle, fh.handle.LinkType())
	for packet := range packetSource.Packets() {
		captureInfo := packet.Metadata().CaptureInfo
		rawPacket := packet.Data()
		for _, forwarder := range fh.forwarders {
			if forwarder.BPF.Matches(captureInfo, rawPacket) {
				forwarder.Handler.Forward(packet)
			}
		}
	}
}

func MakeInterfaceForwarder(handle *pcap.Handle, c *Config) (*InterfaceForwarder, error) {
	forwarders := []*Forwarder{}
	for _, fc := range c.Forwarders {
		f, err := MakeForwarder(handle, &fc)
		if err != nil {
			return nil, err
		}

		forwarders = append(forwarders, f)
	}

	return &InterfaceForwarder{handle: handle, forwarders: forwarders}, nil
}
