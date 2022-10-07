package main

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
		if disabled := v.GetBool("DisableDHCP"); disabled {
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
	port = v.GetString("DestinationPort")
	if port == "" {
		return errors.New("No destination port for dhcp forwarder")
	}

	config.Port = port
	filter = v.GetString("Filter")
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
		return errors.New("No destination host for dhcp forwarder")
	}

	config.Host = host
	port = v.GetString("DNSDestinationPort")
	if port == "" {
		return errors.New("No destination port for dhcp forwarder")
	}

	config.Port = port
	filter = v.GetString("DNSFilter")
	if filter == "" {
		return errors.New("No filter for dhcp forwarder")
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

	device := v.GetString("ListeningDevice")
	if device == "" {
		return nil, errors.New("No device")
	}

	config.Interface = v.GetString("ListeningDevice")
	filters := []string{}
	for _, f := range config.Forwarders {
		filters = append(filters, "("+f.Filter+")")
	}

	config.Filter = strings.Join(filters, " and ")

	return config, nil
}

func (c *Config) SetupPcapForwarding() (*pcap.Handle, []*Forwarder, error) {
	handle, err := pcap.OpenLive(c.Interface, c.SnapLen, true, pcap.BlockForever)
	if err != nil {
		return nil, nil, err
	}

	err = handle.SetBPFFilter(c.Filter)
	if err != nil {
		return nil, nil, err
	}

	forwarders := []*Forwarder{}
	for _, fc := range c.Forwarders {
		f, err := MakeForwarder(handle, &fc)
		if err != nil {
			return nil, nil, err
		}
		forwarders = append(forwarders, f)
	}

	return handle, forwarders, nil
}
