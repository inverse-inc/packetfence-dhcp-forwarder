package forwarder

import (
	"errors"
	"os"
	"strings"

	"github.com/google/gopacket/pcap"
	"github.com/spf13/viper"
)

const defaultSnapLen = 1600

type Config struct {
	Forwarders []ForwarderConfig
	Interface  string
	Filter     string
	SnapLen    int32
}

type ForwarderConfig struct {
	Type   string
	Filter string
	Host   string
	Port   string
}

func initViper(name string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(name) // will match {name}.{toml,json} etc.
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	v.AddConfigPath(pwd)
	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return v, nil
}

func isDisabled(v *viper.Viper, key string) bool {
	if v.IsSet(key) {
		if disabled := v.GetBool(key); disabled {
			return true
		}
	}

	return false
}

func getDHCPConfig(v *viper.Viper, c *Config) error {
	if isDisabled(v, "DisableDHCP") {
		return nil
	}

	config := ForwarderConfig{Type: "dhcp"}
	host := v.GetString("DestinationHost")
	if host == "" {
		return errors.New("No destination host for dhcp forwarder")
	}

	config.Host = host
	port := v.GetString("DestinationPort")
	if port == "" {
		return errors.New("No destination port for dhcp forwarder")
	}

	config.Port = port
	filter := v.GetString("Filter")
	if filter == "" {
		return errors.New("No filter for dhcp forwarder")
	}

	config.Filter = filter
	c.Forwarders = append(c.Forwarders, config)
	return nil
}

func getDNSConfig(v *viper.Viper, c *Config) error {
	if isDisabled(v, "DisableDNS") {
		return nil
	}

	config := ForwarderConfig{Type: "dns"}
	host := v.GetString("DNSDestinationHost")
	if host == "" {
		return errors.New("No destination host for dns forwarder")
	}

	config.Host = host
	port := v.GetString("DNSDestinationPort")
	if port == "" {
		return errors.New("No destination port for dns forwarder")
	}

	config.Port = port
	filter := v.GetString("DNSFilter")
	if filter == "" {
		return errors.New("No filter for dns forwarder")
	}

	config.Filter = filter
	c.Forwarders = append(c.Forwarders, config)
	return nil
}

func GetConfigFromFile(name string) (*Config, error) {
	v, err := initViper(name)
	if err != nil {
		return nil, err
	}

	config := &Config{SnapLen: defaultSnapLen}
	err = getDHCPConfig(v, config)
	if err != nil {
		return nil, err
	}

	err = getDNSConfig(v, config)
	if err != nil {
		return nil, err
	}

	config.Interface = v.GetString("ListeningDevice")
	filters := []string{}
	excludes := []string{}
	for _, f := range config.Forwarders {
		filters = append(filters, "("+f.Filter+")")
		excludes = append(excludes, "(not (dst port "+f.Port+" and dst host "+f.Host+" ))")
	}

	config.Filter = "(" + strings.Join(filters, " or ") + ")"
	if len(excludes) > 0 {
		config.Filter += " and (" + strings.Join(excludes, " and ") + ")"
	}

	return config, nil
}

func (c *Config) getHandle(i string) (*pcap.Handle, error) {
	handle, err := pcap.OpenLive(i, c.SnapLen, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	err = handle.SetBPFFilter(c.Filter)
	if err != nil {
		return nil, err
	}

	return handle, nil
}

func (c *Config) SetupPcapForwarding() ([]*InterfaceForwarder, error) {
	if c.Interface != "" {
		handle, err := c.getHandle(c.Interface)
		if err != nil {
			return nil, err
		}

		iFor, err := MakeInterfaceForwarder(handle, c)
		if err != nil {
			return nil, err
		}
		return []*InterfaceForwarder{iFor}, nil
	}

	interfaces, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	interfaceForwarders := []*InterfaceForwarder{}
	for _, i := range interfaces {
		if len(i.Addresses) == 0 {
			continue
		}
		handle, err := c.getHandle(i.Name)
		if err != nil {
			return nil, err
		}

		iFor, err := MakeInterfaceForwarder(handle, c)
		if err != nil {
			return nil, err
		}

		interfaceForwarders = append(interfaceForwarders, iFor)
	}

	return interfaceForwarders, nil
}
