Testing evironment
==========

Forwarder machine: 
Core 2 duo 2.4ghz 


We injected 1720 pps of DHCP traffic for 30 seconds without any packet being dropped and the forwarding machine having a cpu never being more impacted then 2%-5%, including system time processing the packets.


Tests
==========
Data source:
capinfos dhcp_client.pcap 
File name:           dhcp_client.pcap
File type:           Wireshark/tcpdump/... - pcap
File encapsulation:  Ethernet
File timestamp precision:  microseconds (6)
Packet size limit:   file hdr: 65535 bytes
Number of packets:   50 k
File size:           18 MB
Data size:           17 MB
Capture duration:    2904.258445 seconds
First packet time:   2016-09-28 19:59:25.655624
Last packet time:    2016-09-28 20:47:49.914069
Data byte rate:      5942 bytes/s
Data bit rate:       47 kbps
Average packet size: 345.18 bytes
Average packet rate: 17 packets/s
SHA1:                d31e4961ebaf481542023a075bc589e3fb167f18
RIPEMD160:           ee922a9b28a9d676eee9166dc24783d764709075
MD5:                 445b5d7eaaca1442aff3d85615e05980
Strict time order:   True
Number of interfaces in file: 1
Interface #0 info:
                     Name = UNKNOWN
                     Description = NONE
                     Encapsulation = Ethernet (1/1 - ether)
                     Speed = 0
                     Capture length = 65535
                     FCS length = -1
                     Time precision = microseconds (6)
                     Time ticks per second = 1000000
                     Time resolution = 0x06
                     Filter string = NONE
                     Operating system = UNKNOWN
                     Comment = NONE
                     BPF filter length = 0
                     Number of stat entries = 0
                     Number of packets = 50000


Results:
We injected 2904 seconds worth of traffic in 30 29.06 seconds (x100 speed).

sudo tcpreplay -K -l 1 -x100 -i enp0s25  dhcp_client.pcap 
sending out enp0s25 
processing file: dhcp_client.pcap
Actual: 50000 packets (17258788 bytes) sent in 29.06 seconds.		Rated: 593901.9 bps, 4.53 Mbps, 1720.58 pps
Statistics for network device: enp0s25
	Attempted packets:         50000
	Successful packets:        50000
	Failed packets:            0
	Retried packets (ENOBUFS): 0
	Retried packets (EAGAIN):  0



Receiving side (same machine):
dudo tcpdump -nneti enp0s25 -s 0 -w Desktop/test_fullspeed_cached_1loop_strike5.pcap port 6678 
tcpdump: listening on enp0s25, link-type EN10MB (Ethernet), capture size 262144 bytes
^C9193 packets captured
9193 packets received by filter

sudo tcpdump -nneti enp0s25 -s 0 -w Desktop/test_30x_cached_loop.pcap port 6678 
tcpdump: listening on enp0s25, link-type EN10MB (Ethernet), capture size 262144 bytes
^C9193 packets captured
9193 packets received by filter
0 packets dropped by kernel

There is actually 9193 packets corresponding to the bpf in the accelerated replay:
tcpdump -r dhcp_client.pcap 'udp and port 67 and ((udp[250:1] = 0x3) or (udp[250:1] = 0x5))' -w tcpdump_output_with_bpf.pcap
reading from file dhcp_client.pcap, link-type EN10MB (Ethernet)
ubuntu@ubuntu:~/Desktop$ capinfos tcpdump_output_with_bpf.pcap 
File name:           tcpdump_output_with_bpf.pcap
File type:           Wireshark/tcpdump/... - pcap
File encapsulation:  Ethernet
File timestamp precision:  microseconds (6)
Packet size limit:   file hdr: 65535 bytes
Number of packets:   9193
File size:           3407 kB
Data size:           3260 kB
Capture duration:    2903.705940 seconds
First packet time:   2016-09-28 19:59:25.769034
Last packet time:    2016-09-28 20:47:49.474974
Data byte rate:      1122 bytes/s
Data bit rate:       8983 bits/s
Average packet size: 354.68 bytes
Average packet rate: 3 packets/s
SHA1:                3bf97366c9a54bd91545e78384c1eb8fb81e098e
RIPEMD160:           eb17a02c5d61cfa8eb7f9ab49d73d8f45b53e166
MD5:                 863b9004b2dc9887b188b00a5317a3b9
Strict time order:   True
Number of interfaces in file: 1
Interface #0 info:
                     Name = UNKNOWN
                     Description = NONE
                     Encapsulation = Ethernet (1/1 - ether)
                     Speed = 0
                     Capture length = 65535
                     FCS length = -1
                     Time precision = microseconds (6)
                     Time ticks per second = 1000000
                     Time resolution = 0x06
                     Filter string = NONE
                     Operating system = UNKNOWN
                     Comment = NONE
                     BPF filter length = 0
                     Number of stat entries = 0
                     Number of packets = 9193


This is to be redone with real captured data from one of our biggest client's registration vlan.
