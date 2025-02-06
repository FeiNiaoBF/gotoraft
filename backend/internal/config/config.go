package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// 配置文件
// 适用于整个项目的配置
// 包括raft的配置，http的配置，websocket的配置

// Config 用于配置整个服务
type Config struct {
	// HTTP 服务相关
	HTTPPort        int    `yaml:"http_port"`        // HTTP 端口
	WebSocketPath   string `yaml:"websocket_path"`   // WebSocket 端点
	EnableWebSocket bool   `yaml:"enable_websocket"` // 是否开启 WebSocket

	// Raft 相关
	RaftNodeCount     int           `yaml:"raft_node_count"`     // Raft 节点数量
	RaftElectionFixed time.Duration `yaml:"raft_election_fixed"` // 固定的选举超时（可使用随机值）
	RaftHeartbeat     time.Duration `yaml:"raft_heartbeat"`      // 心跳间隔

	// 其他全局配置
	DebugMode bool `yaml:"debug_mode"` // 是否开启调试模式
}

// 默认配置
var DefaultConfig = &Config{
	HTTPPort:          8080,
	WebSocketPath:     "/ws",
	EnableWebSocket:   true,
	RaftNodeCount:     3,
	RaftElectionFixed: 1500 * time.Millisecond, // 1.5s
	RaftHeartbeat:     100 * time.Millisecond,  // 0.1s
	DebugMode:         true,
}

// LoadConfig 先尝试解析命令行参数，再尝试从环境变量/文件中加载
func LoadConfig() (*Config, error) {
	// 解析命令行参数
	var (
		configPath = flag.String("config", "", "Path to a YAML config file.")
		nodeCount  = flag.Int("nodes", 3, "Number of Raft nodes.")
		debugMode  = flag.Bool("debug", false, "Enable debug mode.")
	)
	flag.Parse()

	// 设置默认值
	var cfg *Config
	if flag.NArg() == 0 {
		cfg = DefaultConfig
	} else {
		cfg = &Config{
			HTTPPort:          8080,
			WebSocketPath:     "/ws",
			EnableWebSocket:   true,
			RaftNodeCount:     *nodeCount,
			RaftElectionFixed: 1500 * time.Millisecond, // 1.5s
			RaftHeartbeat:     100 * time.Millisecond,  // 0.1s
			DebugMode:         *debugMode,
		}
	}

	// 如果指定了配置文件，则从文件中加载并覆盖默认值
	if *configPath != "" {
		fileCfg, err := loadFromYAML(*configPath)
		if err != nil {
			return nil, fmt.Errorf("loadFromYAML error: %v", err)
		}
		mergeConfig(cfg, fileCfg)
	}

	// 校验最终配置
	if cfg.RaftNodeCount <= 0 {
		return nil, fmt.Errorf("RaftNodeCount must be > 0")
	} 

	return cfg, nil
}

// loadFromYAML 从 YAML 文件中加载配置
func loadFromYAML(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fileCfg Config
	if err := yaml.Unmarshal(data, &fileCfg); err != nil {
		return nil, err
	}
	return &fileCfg, nil
}

// mergeConfig 将 fileCfg 中非零值/非空值合并到 cfg 中
func mergeConfig(cfg, fileCfg *Config) {
	if fileCfg.HTTPPort != 0 {
		cfg.HTTPPort = fileCfg.HTTPPort
	}
	if fileCfg.WebSocketPath != "" {
		cfg.WebSocketPath = fileCfg.WebSocketPath
	}
	cfg.EnableWebSocket = fileCfg.EnableWebSocket

	if fileCfg.RaftNodeCount != 0 {
		cfg.RaftNodeCount = fileCfg.RaftNodeCount
	}
	if fileCfg.RaftElectionFixed != 0 {
		cfg.RaftElectionFixed = fileCfg.RaftElectionFixed
	}
	if fileCfg.RaftHeartbeat != 0 {
		cfg.RaftHeartbeat = fileCfg.RaftHeartbeat
	}
	cfg.DebugMode = fileCfg.DebugMode
}
