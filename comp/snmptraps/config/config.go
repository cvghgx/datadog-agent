// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2020-present Datadog, Inc.

package config

import (
	"fmt"
	"hash/fnv"

	"github.com/gosnmp/gosnmp"

	"github.com/DataDog/datadog-agent/comp/core/config"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	"github.com/DataDog/datadog-agent/comp/snmptraps/snmplog"
	"github.com/DataDog/datadog-agent/pkg/snmp/gosnmplib"
	"github.com/DataDog/datadog-agent/pkg/snmp/utils"
)

const (
	defaultPort        = uint16(9162) // Standard UDP port for traps.
	defaultStopTimeout = 5
	packetsChanSize    = 100
)

// UserV3 contains the definition of one SNMPv3 user with its username and its auth
// parameters.
type UserV3 struct {
	Username       string `mapstructure:"user" yaml:"user"`
	UsernameLegacy string `mapstructure:"username" yaml:"username"`
	AuthKey        string `mapstructure:"authKey" yaml:"authKey"`
	AuthProtocol   string `mapstructure:"authProtocol" yaml:"authProtocol"`
	PrivKey        string `mapstructure:"privKey" yaml:"privKey"`
	PrivProtocol   string `mapstructure:"privProtocol" yaml:"privProtocol"`
}

// TrapsConfig contains configuration for SNMP trap listeners.
// YAML field tags provided for test marshalling purposes.
type TrapsConfig struct {
	Enabled               bool     `mapstructure:"enabled" yaml:"enabled"`
	Port                  uint16   `mapstructure:"port" yaml:"port"`
	Users                 []UserV3 `mapstructure:"users" yaml:"users"`
	CommunityStrings      []string `mapstructure:"community_strings" yaml:"community_strings"`
	BindHost              string   `mapstructure:"bind_host" yaml:"bind_host"`
	StopTimeout           int      `mapstructure:"stop_timeout" yaml:"stop_timeout"`
	Namespace             string   `mapstructure:"namespace" yaml:"namespace"`
	authoritativeEngineID string   `mapstructure:"-" yaml:"-"`
}

// ReadConfig builds the traps configuration from the Agent configuration.
func ReadConfig(host string, conf config.Component) (*TrapsConfig, error) {
	var c = &TrapsConfig{}
	err := conf.UnmarshalKey("network_devices.snmp_traps", &c)
	if err != nil {
		return nil, err
	}
	if err := c.SetDefaults(host, conf.GetString("network_devices.namespace")); err != nil {
		return c, err
	}
	return c, nil
}

// SetDefaults sets all unset values to default values, and returns an error
// if any fields are invalid.
func (c *TrapsConfig) SetDefaults(host string, namespace string) error {
	// Set defaults.
	if c.Port == 0 {
		c.Port = defaultPort
	}
	if c.BindHost == "" {
		// Default to global bind_host option.
		c.BindHost = "0.0.0.0"
	}
	if c.StopTimeout == 0 {
		c.StopTimeout = defaultStopTimeout
	}
	if host == "" {
		// Make sure to have at least some unique bytes for the authoritative engineID.
		// Unlikely to happen since the agent cannot start without a hostname
		host = "unknown-datadog-agent"
	}
	h := fnv.New128()
	h.Write([]byte(host))
	// First byte is always 0x80
	// Next four bytes are the Private Enterprise Number (set to an invalid value here)
	// The next 16 bytes are the hash of the agent hostname
	engineID := h.Sum([]byte{0x80, 0xff, 0xff, 0xff, 0xff})
	c.authoritativeEngineID = string(engineID)

	if c.Namespace == "" {
		c.Namespace = namespace
	}
	var err error
	c.Namespace, err = utils.NormalizeNamespace(c.Namespace)
	if err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	return nil
}

// Addr returns the host:port address to listen on.
func (c *TrapsConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.BindHost, c.Port)
}

// BuildSNMPParams returns a valid GoSNMP params structure from configuration.
func (c *TrapsConfig) BuildSNMPParams(logger log.Component) (*gosnmp.GoSNMP, error) {
	var snmpLogger gosnmp.Logger
	if logger != nil {
		snmpLogger = gosnmp.NewLogger(snmplog.New(logger))
	}
	if len(c.Users) == 0 {
		return &gosnmp.GoSNMP{
			Port:      c.Port,
			Transport: "udp",
			Version:   gosnmp.Version2c, // No user configured, let's use Version2 which is enough and doesn't require setting up fake security data.
			Logger:    snmpLogger,
		}, nil
	}

	// Set up user security params table from config
	usmTable := gosnmp.NewSnmpV3SecurityParametersTable(snmpLogger)
	for _, user := range c.Users {
		// Backward compatibility
		if user.Username == "" {
			user.Username = user.UsernameLegacy
		}

		authProtocol, err := gosnmplib.GetAuthProtocol(user.AuthProtocol)
		if err != nil {
			return nil, err
		}
		privProtocol, err := gosnmplib.GetPrivProtocol(user.PrivProtocol)
		if err != nil {
			return nil, err
		}
		err = usmTable.Add(user.Username, &gosnmp.UsmSecurityParameters{
			UserName:                 user.Username,
			AuthenticationProtocol:   authProtocol,
			AuthenticationPassphrase: user.AuthKey,
			PrivacyProtocol:          privProtocol,
			PrivacyPassphrase:        user.PrivKey,
		})
		if err != nil {
			return nil, err
		}
	}

	return &gosnmp.GoSNMP{
		Port:                        c.Port,
		Transport:                   "udp",
		Version:                     gosnmp.Version3, // Always using version3 for traps, only option that works with all SNMP versions simultaneously
		SecurityModel:               gosnmp.UserSecurityModel,
		SecurityParameters:          &gosnmp.UsmSecurityParameters{AuthoritativeEngineID: c.authoritativeEngineID},
		TrapSecurityParametersTable: usmTable,
		Logger:                      snmpLogger,
	}, nil
}

// GetPacketChannelSize returns the default size for the packets channel
func (c *TrapsConfig) GetPacketChannelSize() int {
	return packetsChanSize
}

// IsEnabled returns whether SNMP trap collection is enabled in the Agent configuration.
func IsEnabled(conf config.Component) bool {
	return conf.GetBool("network_devices.snmp_traps.enabled")
}
