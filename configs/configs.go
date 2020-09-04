// Package configs contains app configuration
package configs

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

var (
	configs Configs
	once    sync.Once
)

// Configs consists all configuration
type Configs struct {
	Debug  bool
	Server ServerConfig
	Mysql  MySQLConfig
}

// ServerConfig consists server configuration
type ServerConfig struct {
	Address string
}

// MySQLConfig consists MySQL database configuration
type MySQLConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

// GetConfig return Configs object read from config.json file
func GetConfig() *Configs {
	once.Do(func() {
		conf := viper.New()
		conf.SetConfigFile("./configs/config.json")

		err := conf.ReadInConfig()
		if err != nil {
			log.Fatalf("failed to read config file: %s", err)
		}

		if err := conf.Unmarshal(&configs); err != nil {
			log.Fatalf("failed to unmarshal config: %s", err)
		}
	})

	return &configs
}
