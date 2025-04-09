package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// DBConfig holds database-specific configuration
type DBConfig struct {
	Path             string `mapstructure:"db_path"`
	MigrationsPath   string `mapstructure:"migrations_path"`
	InitialMigration string `mapstructure:"initial_migration"`
}

// SSHConfig holds SSH Git server configuration
type SSHConfig struct {
	Address string `mapstructure:"ssh_git_address"`
	HostKey string `mapstructure:"ssh_git_hostkey"`
}

// Configuration holds all module-specific configurations
type Configuration struct {
	DB  DBConfig  `mapstructure:"db"`
	SSH SSHConfig `mapstructure:"ssh"`
}

// Load initializes the configuration from environment variables and config files
func Load() (*Configuration, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("db.path", "data/depgit.db")
	v.SetDefault("db.migrations_path", "")
	v.SetDefault("db.initial_migration", "")
	v.SetDefault("ssh.address", "0.0.0.0:2222")
	v.SetDefault("ssh.hostkey", "")

	// Enable environment variable support with nested key support
	v.SetEnvPrefix("DEPGIT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Explicitly bind environment variables to their config keys
	v.BindEnv("db.path", "DEPGIT_DB_PATH")
	v.BindEnv("db.migrations_path", "DEPGIT_MIGRATIONS_PATH")
	v.BindEnv("db.initial_migration", "DEPGIT_INITIAL_MIGRATION")
	v.BindEnv("ssh.address", "DEPGIT_SSH_GIT_ADDRESS")
	v.BindEnv("ssh.hostkey", "DEPGIT_SSH_GIT_HOSTKEY")

	// Map old env vars to new structure for backward compatibility
	if path := os.Getenv("DEPGIT_DB_PATH"); path != "" {
		v.Set("db.path", path)
	}
	if path := os.Getenv("DEPGIT_MIGRATIONS_PATH"); path != "" {
		v.Set("db.migrations_path", path)
	}
	if migration := os.Getenv("DEPGIT_INITIAL_MIGRATION"); migration != "" {
		v.Set("db.initial_migration", migration)
	}
	if addr := os.Getenv("DEPGIT_SSH_GIT_ADDRESS"); addr != "" {
		v.Set("ssh.address", addr)
	}
	if key := os.Getenv("DEPGIT_SSH_GIT_HOSTKEY"); key != "" {
		v.Set("ssh.hostkey", key)
	}

	// Try to read config file if it exists
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	// Read the config file, but don't error if it doesn't exist
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := &Configuration{}
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Create database directory if necessary
	dbDir := filepath.Dir(config.DB.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating database directory: %w", err)
	}

	return config, nil
}

// GetDatabasePath returns the configured database path
func (c *Configuration) GetDatabasePath() string {
	return c.DB.Path
}

// GetMigrationsPath returns the configured migrations path
func (c *Configuration) GetMigrationsPath() string {
	return c.DB.MigrationsPath
}

// IsTrue checks if an environment variable represents a true value
func IsTrue(envVar string) bool {
	val := strings.ToLower(os.Getenv(envVar))
	return val == "true" || val == "yes" || val == "1"
}
