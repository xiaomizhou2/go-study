package config

import (
	"fmt"

	"github.com/example/user-service/logger"
	"github.com/spf13/viper"
)

// Config 应用配置
// 对比 Java: @ConfigurationProperties(prefix = "app")
//
// Viper 会自动把 YAML 的键名映射到结构体字段：
//   server.port    → ServerConfig.Port
//   database.dsn   → DatabaseConfig.DSN
// 注意：结构体字段名用 mapstructure tag 告诉 Viper 怎么映射
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      logger.Config  `mapstructure:"log"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port string `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`
	DSN          string `mapstructure:"dsn"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// Load 加载配置文件
// 对比 Java: Spring 自动读取 application.yml 并注入到 @ConfigurationProperties
// Go 里我们手动读，但只有几行代码
//
// configPath 参数让调用方决定配置文件位置，默认传 "./config.yaml"
// 对比 Java: spring.config.location=classpath:/config.yaml
func Load(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	// 支持环境变量覆盖（对比 Java: spring.config.import=optional:env）
	// 比如 APP_SERVER_PORT=:9090 会覆盖 config.yaml 里的 server.port
	v.SetEnvPrefix("APP")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	return &cfg, nil
}
