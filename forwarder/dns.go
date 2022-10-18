package forwarder

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const PFForwardSrcDst layers.DNSOptionCode = 65245

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

func addSrcDstRR(dns *layers.DNS, code layers.DNSOptionCode, src, dst net.IP) {
	if dns == nil {
		return
	}

	var rr *layers.DNSResourceRecord = nil

	// find the last resource record
	for i := len(dns.Additionals) - 1; i >= 0; i++ {
		if dns.Additionals[i].Type == layers.DNSTypeOPT {
			rr = &dns.Additionals[i]
			break
		}
	}

	if rr == nil {
		dns.Additionals = append(
			dns.Additionals,
			layers.DNSResourceRecord{
				Type: layers.DNSTypeOPT,
			},
		)

		rr = &dns.Additionals[len(dns.Additionals)-1]
	}

	rr.OPT = append(
		rr.OPT,
		layers.DNSOPT{
			Code: code,
			Data: append(src, dst...),
		},
	)
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

	addSrcDstRR(&dns, PFForwardSrcDst, ipv4.SrcIP, ipv4.DstIP)
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true}
	if err := dns.SerializeTo(buf, opts); err != nil {
		return err
	}

	h.conn.WriteToUDP(buf.Bytes(), h.addr)
	return nil
}
