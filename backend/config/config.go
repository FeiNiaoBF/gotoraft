// config.go 是配置文件，用来配置整个系统的参数和上下文
// 使用 viper （github.com/spf13/viper）来加载和管理配置
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config 用于存储和管理配置
// TODO: 可能需要加入更多Server的配置
// `mapstructure` 是用来默认支持自动的蛇形命名（snake_case）到驼峰命名（CamelCase）的转换
type Config struct {
	Server *ServerConfig `mapstructure:"server"`
	Log    *LogConfig    `mapstructure:"log"`
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

var (
	// AppConfig 全局配置实例
	AppConfig Config
)

// Init 初始化配置
// config.yml 里面的配置会被加载到 AppConfig 中
func Init() error {
	// 设置配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("未找到配置文件，使用默认配置")
		} else {
			return err
		}
	}

	// 设置默认值
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

	// 解析配置到结构体
	if err := viper.Unmarshal(&AppConfig); err != nil {
		return err
	}

	// 确保日志目录存在
	logDir := filepath.Dir(AppConfig.Log.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

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
