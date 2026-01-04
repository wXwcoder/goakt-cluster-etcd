/*
 * MIT License
 *
 * Copyright (c) 2022-2025 Arsene Tochemey Gandote
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package service

import "github.com/caarlos0/env/v11"

// Config defines the service configuration
type Config struct {
	Port            int      `env:"PORT" envDefault:"50051"`
	ServiceName     string   `env:"SERVICE_NAME"`
	ActorSystemName string   `env:"SYSTEM_NAME"`
	GossipPort      int      `env:"GOSSIP_PORT"`
	PeersPort       int      `env:"PEERS_PORT"`
	RemotingPort    int      `env:"REMOTING_PORT"`
	EtcdEndpoints   []string `env:"ETCD_ENDPOINTS" envSeparator:"," envDefault:"localhost:2379"`
	EtcdDialTimeout int      `env:"ETCD_DIAL_TIMEOUT" envDefault:"5"`
	EtcdTTL         int64    `env:"ETCD_TTL" envDefault:"30"`
	EtcdTimeout     int      `env:"ETCD_TIMEOUT" envDefault:"5"`
}

// GetConfig returns the configuration from environment variables
func GetConfig() (*Config, error) {
	// load the host node configuration from environment variables
	cfg := &Config{}
	opts := env.Options{RequiredIfNoDef: true, UseFieldNameByDefault: false}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetConfigWithOverrides returns configuration with command line overrides
func GetConfigWithOverrides(port int, serviceName, systemName string, gossipPort, peersPort, remotingPort int, etcdEndpoints []string, etcdDialTimeout int, etcdTTL int64, etcdTimeout int) (*Config, error) {
	// load the host node configuration from environment variables
	cfg := &Config{}
	opts := env.Options{RequiredIfNoDef: true, UseFieldNameByDefault: false}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, err
	}

	// Override with command line parameters if provided
	if port > 0 {
		cfg.Port = port
	}
	if serviceName != "" {
		cfg.ServiceName = serviceName
	}
	if systemName != "" {
		cfg.ActorSystemName = systemName
	}
	if gossipPort > 0 {
		cfg.GossipPort = gossipPort
	}
	if peersPort > 0 {
		cfg.PeersPort = peersPort
	}
	if remotingPort > 0 {
		cfg.RemotingPort = remotingPort
	}
	if len(etcdEndpoints) > 0 {
		cfg.EtcdEndpoints = etcdEndpoints
	}
	if etcdDialTimeout > 0 {
		cfg.EtcdDialTimeout = etcdDialTimeout
	}
	if etcdTTL > 0 {
		cfg.EtcdTTL = etcdTTL
	}
	if etcdTimeout > 0 {
		cfg.EtcdTimeout = etcdTimeout
	}

	return cfg, nil
}
