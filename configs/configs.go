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
	Production    bool
	Server        ServerConfig
	Mysql         MySQLConfig
	ElasticSearch ESConfig
}

// ServerConfig consists server configuration
type ServerConfig struct {
	Address      string
	WriteTimeout int
	ReadTimeout  int
	IdleTimeout  int
}

// MySQLConfig consists MySQL database configuration
type MySQLConfig struct {
	Host            string
	Port            string
	User            string
	Pass            string
	Name            string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime int
}

// ESConfig consists ElasticSearch configuration
type ESConfig struct {
	Address  string
	IsAuth   bool
	Username string
	Password string
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
