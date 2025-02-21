// config.go 是配置文件，用来配置整个系统的参数和上下文
// 使用 viper （github.com/spf13/viper）来加载和管理配置
package config

import (
	"gotoraft/internal/raft"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 用于存储和管理配置
type Config struct {
	Server *ServerConfig `mapstructure:"server"`
	Log    *LogConfig    `mapstructure:"log"`
	Raft   *raft.Config  `mapstructure:"raft"`
	NetSim *NetSimConfig `mapstructure:"net_sim"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
	TimeFormat string `mapstructure:"time_format"`
}

// NetSimConfig 网络模拟配置
type NetSimConfig struct {
	// 最小网络延迟（毫秒）
	NetworkLatencyMin int `mapstructure:"network_latency_min"`

	// 最大网络延迟（毫秒）
	NetworkLatencyMax int `mapstructure:"network_latency_max"`

	// 丢包率 (0-1)
	PacketLossRate float64 `mapstructure:"packet_loss_rate"`
}

var (
	// AppConfig 全局配置实例
	AppConfig Config
)

// Init 初始化配置
func Init() error {
	// 设置配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 添加配置文件路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.AddConfigPath("$HOME/.gotoraft")
	viper.AddConfigPath("/etc/gotoraft")

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在时创建默认配置文件
			if err := createDefaultConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	return nil
}

// setDefaults 设置默认值
func setDefaults() {
	// 服务器默认配置
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 8080)

	// 日志默认配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.filename", "gotoraft.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.max_backups", 10)
	viper.SetDefault("log.compress", true)
	viper.SetDefault("log.time_format", "2006-01-02 15:04:05")

	// Raft默认配置
	defaultRaftConfig := raft.DefaultConfig()
	viper.SetDefault("raft", defaultRaftConfig)

	// 网络模拟默认配置
	viper.SetDefault("net_sim.network_latency_min", 10)
	viper.SetDefault("net_sim.network_latency_max", 100)
	viper.SetDefault("net_sim.packet_loss_rate", 0.1)
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig() error {
	configDir := "config"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	if err := viper.SafeWriteConfigAs(configPath); err != nil {
		return err
	}

	log.Printf("Created default config file at: %s", configPath)
	return nil
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return &AppConfig
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	if AppConfig.Server == nil {
		return nil
	}
	return AppConfig.Server
}

// GetLogConfig 获取日志配置
func GetLogConfig() *LogConfig {
	if AppConfig.Log == nil {
		return nil
	}
	return AppConfig.Log
}

// GetRaftConfig 获取Raft配置
func GetRaftConfig() *raft.Config {
	if AppConfig.Raft == nil {
		return raft.DefaultConfig()
	}
	return AppConfig.Raft
}

// GetNetSimConfig 获取网络模拟配置
func GetNetSimConfig() *NetSimConfig {
	if AppConfig.NetSim == nil {
		return nil
	}
	return AppConfig.NetSim
}
