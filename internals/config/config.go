package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type HTTPS struct {
	Address         string        `yaml:"address" env:"HTTP_ADDRESS"`
	Port            int           `yaml:"port" env:"HTTP_PORT"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout     time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT"`
}

type Database struct {
	Host            string        `yaml:"host" env:"DB_HOST"`
	Port            int           `yaml:"port" env:"DB_PORT"`
	User            string        `yaml:"user" env:"DB_USER"`
	Password        string        `yaml:"password" env:"DB_PASSWORD"`
	Name            string        `yaml:"name" env:"DB_NAME"`
	SSLMode         string        `yaml:"ssl_mode" env:"DB_SSL_MODE"`
	MaxOpenConns    int           `yaml:"max_open_conns" env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env:"DB_MAX_IDLE_CONNS"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime" env:"DB_CONN_MAX_LIFETIME"`
}

type Logger struct {
	Level      string `yaml:"level" env:"LOG_LEVEL"`
	Format     string `yaml:"format" env:"LOG_FORMAT"` // json or text
	OutputPath string `yaml:"output_path" env:"LOG_OUTPUT_PATH"`
}

type Config struct {
	Env         string   `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath string   `yaml:"storage_path" env:"STORAGE_PATH"`
	HTTPS       HTTPS    `yaml:"http_server"`
	Database    Database `yaml:"database"`
	Logger      Logger   `yaml:"logger"`
}

// MustLoad loads configuration from file and panics on error
func MustLoad(configPath string) *Config {
	cfg, err := Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

// Load reads and parses the configuration file
func Load(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cfg.SetDefaults()

	return &cfg, nil
}

// Validate checks if required fields are set
func (c *Config) Validate() error {
	if c.Env == "" {
		return fmt.Errorf("env is required")
	}
	return nil
}

// SetDefaults sets default values for optional fields
func (c *Config) SetDefaults() {
	if c.HTTPS.Address == "" {
		c.HTTPS.Address = "localhost"
	}
	if c.HTTPS.Port == 0 {
		c.HTTPS.Port = 8080
	}
	if c.HTTPS.ReadTimeout == 0 {
		c.HTTPS.ReadTimeout = 10 * time.Second
	}
	if c.HTTPS.WriteTimeout == 0 {
		c.HTTPS.WriteTimeout = 10 * time.Second
	}
	if c.HTTPS.IdleTimeout == 0 {
		c.HTTPS.IdleTimeout = 60 * time.Second
	}
	if c.HTTPS.ShutdownTimeout == 0 {
		c.HTTPS.ShutdownTimeout = 30 * time.Second
	}

	if c.Logger.Level == "" {
		c.Logger.Level = "info"
	}
	if c.Logger.Format == "" {
		c.Logger.Format = "json"
	}

	if c.Database.SSLMode == "" {
		c.Database.SSLMode = "disable"
	}
	if c.Database.MaxOpenConns == 0 {
		c.Database.MaxOpenConns = 25
	}
	if c.Database.MaxIdleConns == 0 {
		c.Database.MaxIdleConns = 5
	}
	if c.Database.ConnMaxLifetime == 0 {
		c.Database.ConnMaxLifetime = 5 * time.Minute
	}
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.HTTPS.Address, c.HTTPS.Port)
}
