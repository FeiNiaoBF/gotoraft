// config.go 是配置文件，用来配置整个系统的参数和上下文
// 使用 viper （github.com/spf13/viper）来加载和管理配置
package config

import (
	"github.com/spf13/viper"
)

// Config 用于存储和管理配置
type Config struct {
	Server *ServerConfig `mapstructure:"server"`
	Log    *LogConfig    `mapstructure:"log"`
	Raft   *RaftConfig   `mapstructure:"raft"`
	NodeID string        `mapstructure:"node_id"`
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

// RaftConfig Raft配置
type RaftConfig struct {
	DataDir   string `mapstructure:"data_dir"`  // Raft数据存储目录
	Host      string `mapstructure:"host"`      // Raft节点监听地址
	Port      int    `mapstructure:"port"`      // Raft节点监听端口
	Bootstrap bool   `mapstructure:"bootstrap"` // 是否为初始化节点
	JoinAddr  string `mapstructure:"join_addr"` // 加入集群的地址
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
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件不存在时创建默认配置
			if err := createDefaultConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 将配置解析到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	return nil
}

// setDefaults 设置默认值
func setDefaults() {
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "console")
	viper.SetDefault("log.output", "both")
	viper.SetDefault("log.filename", "logs/server.log")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_age", 7)
	viper.SetDefault("log.max_backups", 10)
	viper.SetDefault("log.compress", true)
	viper.SetDefault("log.time_format", "2006-01-02 15:04:05")

	viper.SetDefault("raft.data_dir", "data/raft")
	viper.SetDefault("raft.host", "0.0.0.0")
	viper.SetDefault("raft.port", 8081)
	viper.SetDefault("raft.bootstrap", false)
	viper.SetDefault("raft.join_addr", "")

	viper.SetDefault("node_id", "node1")
}

// createDefaultConfig 创建默认配置文件
func createDefaultConfig() error {
	setDefaults()
	return viper.SafeWriteConfig()
}

// GetConfig 获取配置实例
func GetConfig() *Config {
	return &AppConfig
}

// GetServerConfig 获取服务器配置
func GetServerConfig() *ServerConfig {
	return AppConfig.Server
}

// GetLogConfig 获取日志配置
func GetLogConfig() *LogConfig {
	return AppConfig.Log
}

// GetRaftConfig 获取Raft配置
func GetRaftConfig() *RaftConfig {
	return AppConfig.Raft
}
