// Package config provides configuration management for the Canopy node.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/canopy-network/canopy/lib/crypto"
)

const (
	DefaultRPCPort     = 50832
	DefaultP2PPort     = 9001
	DefaultLogLevel    = "info"
	DefaultDataDir     = ".canopy"
	DefaultConfigFile  = "config.json"
)

// Config holds all configuration parameters for a Canopy node.
type Config struct {
	// Network identity
	ChainID   uint64 `json:"chain_id"`
	NetworkID string `json:"network_id"`

	// Storage
	DataDirPath string `json:"data_dir"`

	// RPC server
	RPCUrl  string `json:"rpc_url"`
	RPCPort int    `json:"rpc_port"`

	// P2P networking
	P2PPort    int      `json:"p2p_port"`
	BootPeers  []string `json:"boot_peers"`
	MaxPeers   int      `json:"max_peers"`

	// Consensus
	TimeoutPropose  int `json:"timeout_propose_ms"`
	TimeoutVote     int `json:"timeout_vote_ms"`
	TimeoutCommit   int `json:"timeout_commit_ms"`

	// Logging
	LogLevel string `json:"log_level"`

	// Validator
	ValidatorKey string `json:"validator_key,omitempty"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		ChainID:         1,
		NetworkID:       "canopy-mainnet",
		DataDirPath:     filepath.Join(home, DefaultDataDir),
		RPCUrl:          "0.0.0.0",
		RPCPort:         DefaultRPCPort,
		P2PPort:         DefaultP2PPort,
		BootPeers:       []string{},
		MaxPeers:        30,
		TimeoutPropose:  3000,
		TimeoutVote:     2000,
		TimeoutCommit:   1000,
		LogLevel:        DefaultLogLevel,
	}
}

// LoadConfig reads and unmarshals a JSON config file from the given path.
// If the file does not exist, the default config is returned.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config as formatted JSON to the given path,
// creating any necessary parent directories.
func (c *Config) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// ConfigFilePath returns the canonical path to the config file
// inside the node's data directory.
func (c *Config) ConfigFilePath() string {
	return filepath.Join(c.DataDirPath, DefaultConfigFile)
}

// ValidatorKeyPath returns the path to the validator private key file.
func (c *Config) ValidatorKeyPath() string {
	return filepath.Join(c.DataDirPath, crypto.PrivKeyFileName)
}
