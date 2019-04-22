package config

import (
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

// Config defines the main application config file
type Config struct {
	Mode          string
	Address       string
	SessionSecret string
	Database      DatabaseConfig
	Roles         map[string][]string
}

// DatabaseConfig defines the main application database
type DatabaseConfig struct {
	Username string
	Password string
	Database string
}

// LoadConfigFile loads the given config file
func LoadConfigFile(path string) (*Config, error) {
	cfg := Config{}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ConfigUseContext create a gin config context
func ConfigUseContext(cfg *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", cfg)
		c.Next()
	}
}
