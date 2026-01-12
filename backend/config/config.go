package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Driver string
	Source string
}

// LoadConfig 读取配置文件
func LoadConfig() (*Config, error) {
	viper.AddConfigPath(".")      // 寻找路径：当前目录
	viper.SetConfigName("config") // 寻找文件：config
	viper.SetConfigType("yaml")   // 文件类型：yaml

	// 读取环境变量，支持自动覆盖 (例如 MYBLOG_DATABASE_SOURCE)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MYBLOG")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
