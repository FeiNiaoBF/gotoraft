package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	mu     sync.RWMutex
	Config struct {
		HTTPPort       int    `yaml:"http_port"`
		WebsocketPath  string `yaml:"websocket_path"`
		EnableWebsocket bool  `yaml:"enable_websocket"`

		Raft struct {
			NodeCount           int     `yaml:"node_count"`
			InitialLeader       string  `yaml:"initial_leader"`
			ElectionTimeoutMin  int     `yaml:"election_timeout_min"`
			ElectionTimeoutMax  int     `yaml:"election_timeout_max"`
			HeartbeatInterval   int     `yaml:"heartbeat_interval"`
			MaxLogEntries       int     `yaml:"max_log_entries"`
			ReplicationBatchSize int    `yaml:"replication_batch_size"`
			AnimationSpeed      float64 `yaml:"animation_speed"`
			EnableNetworkDelay  bool    `yaml:"enable_network_delay"`
			NetworkDelayMin     int     `yaml:"network_delay_min"`
			NetworkDelayMax     int     `yaml:"network_delay_max"`
		} `yaml:"raft"`

		KVStore struct {
			InitialData     map[string]string `yaml:"initial_data"`
			MaxKeyLength    int              `yaml:"max_key_length"`
			MaxValueLength  int              `yaml:"max_value_length"`
			MaxEntries      int              `yaml:"max_entries"`
			ReadDelay       int              `yaml:"read_delay"`
			WriteDelay      int              `yaml:"write_delay"`
		} `yaml:"kvstore"`

		WebSocket struct {
			PingInterval     int `yaml:"ping_interval"`
			MaxMessageSize   int `yaml:"max_message_size"`
			WriteBufferSize int `yaml:"write_buffer_size"`
			ReadBufferSize  int `yaml:"read_buffer_size"`
		} `yaml:"websocket"`
	} `yaml:"config"`
}

var (
	config *Config
	once   sync.Once
)

// LoadConfig loads the configuration from the specified file
func LoadConfig() (*Config, error) {
	once.Do(func() {
		config = &Config{}
	})

	data, err := os.ReadFile("configs/env.yaml")
	if err != nil {
		return nil, err
	}

	config.mu.Lock()
	defer config.mu.Unlock()

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	if config == nil {
		LoadConfig()
	}
	return config
}

// UpdateConfig updates the configuration with new values
func (c *Config) UpdateConfig(newConfig Config) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Update configuration values
	c.Config = newConfig.Config

	// Save to file
	data, err := yaml.Marshal(c)
	if err != nil {
		return
	}

	os.WriteFile("configs/env.yaml", data, 0644)
}
