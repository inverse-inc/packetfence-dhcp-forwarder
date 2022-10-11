package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfigLoad(t *testing.T) {
	c, err := GetConfigFromFile("test_config")
	if err != nil {
		t.Fatal("load config")
	}

	if diff := cmp.Diff(
		c,
		&Config{
			Interface: "eth0",
			SnapLen:   1600,
			Filter:    "((udp and port 68 and ((udp[250:1] = 0x3) or (udp[250:1] = 0x5))) or (udp and port 53)) and ((not (dst port 767 and dst host 1.1.1.1 )) and (not (dst port 753 and dst host 1.1.1.1 )))",
			Forwarders: []ForwarderConfig{
				{
					Type:   "dhcp",
					Port:   "767",
					Host:   "1.1.1.1",
					Filter: "udp and port 68 and ((udp[250:1] = 0x3) or (udp[250:1] = 0x5))",
				},
				{
					Type:   "dns",
					Port:   "753",
					Host:   "1.1.1.1",
					Filter: "udp and port 53",
				},
			},
		},
	); diff != "" {
		t.Fatalf("Wrong config: %s", diff)
	}
}
