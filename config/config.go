package config

import (
	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	ServerPort string `mapstructure:"PORT"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Set default values
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("DB_USER", "user")
	viper.SetDefault("DB_PASSWORD", "password")
	viper.SetDefault("DB_NAME", "report_db")
	viper.SetDefault("PORT", "3000")

	// Tell Viper to look for environment variables
	viper.AutomaticEnv() 

	// Explicitly bind environment variables to Viper keys.
	// These keys should match the mapstructure tags or struct field names
	// if tags are not used, for Unmarshal to work correctly.
	// Since we use mapstructure tags, we bind to the tag names.
	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("PORT")

	var config Config
	// Unmarshal the configuration into the Config struct
	// Viper will use the bound environment variables and map them
	// to the struct fields using the mapstructure tags.
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
