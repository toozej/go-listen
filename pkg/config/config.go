package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"`
	Spotify  SpotifyConfig  `envPrefix:"SPOTIFY_"`
	Security SecurityConfig `envPrefix:"SECURITY_"`
	Logging  LoggingConfig  `envPrefix:"LOGGING_"`
}

type ServerConfig struct {
	Host string `env:"HOST" envDefault:"127.0.0.1"`
	Port int    `env:"PORT" envDefault:"8080"`
}

type SpotifyConfig struct {
	ClientID     string `env:"CLIENT_ID"`
	ClientSecret string `env:"CLIENT_SECRET"`
	RedirectURL  string `env:"REDIRECT_URL" envDefault:"http://127.0.0.1:8080/callback"`
}

type SecurityConfig struct {
	RateLimit RateLimitConfig `envPrefix:"RATE_LIMIT_"`
}

type RateLimitConfig struct {
	RequestsPerSecond int `env:"REQUESTS_PER_SECOND" envDefault:"10"`
	Burst             int `env:"BURST" envDefault:"20"`
}

type LoggingConfig struct {
	Level      string `env:"LEVEL" envDefault:"info"`
	Format     string `env:"FORMAT" envDefault:"text"`
	Output     string `env:"OUTPUT" envDefault:"stdout"`
	EnableHTTP bool   `env:"ENABLE_HTTP" envDefault:"true"`
}

// Address returns the server address
func (s ServerConfig) Address() string {
	if s.Host == "" {
		s.Host = "127.0.0.1"
	}
	if s.Port == 0 {
		s.Port = 8080
	}
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func GetEnvVars(debug bool) Config {
	// Load .env file if it exists (will not override existing environment variables)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Warning: Error loading .env file: %s\n", err)
		} else {
			fmt.Println("Loaded environment variables from .env file")
		}
	}

	// Parse environment variables into config struct
	var conf Config
	if err := env.Parse(&conf); err != nil {
		fmt.Printf("Error parsing configuration from environment: %s\n", err)
		os.Exit(1)
	}

	// Validate configuration
	if err := validateConfig(&conf); err != nil {
		fmt.Printf("Configuration validation error: %s\n", err)
		fmt.Println("Please check your configuration and try again.")
		os.Exit(1)
	}

	if debug {
		fmt.Printf("Loaded configuration: %#v\n", conf)
	}

	return conf
}

// validateConfig validates the configuration
func validateConfig(conf *Config) error {
	var errors []string

	// Validate server configuration
	if conf.Server.Port < 1 || conf.Server.Port > 65535 {
		errors = append(errors, "server port must be between 1 and 65535")
	}

	// Validate Spotify configuration (warn but don't fail)
	if conf.Spotify.ClientID == "" {
		fmt.Println("Warning: SPOTIFY_CLIENT_ID is not set. The application will not be able to connect to Spotify.")
		fmt.Println("Please set your Spotify credentials to use the application.")
	}
	if conf.Spotify.ClientSecret == "" {
		fmt.Println("Warning: SPOTIFY_CLIENT_SECRET is not set. The application will not be able to connect to Spotify.")
	}

	// Validate security configuration
	if conf.Security.RateLimit.RequestsPerSecond < 1 {
		errors = append(errors, "rate limit requests per second must be at least 1")
	}
	if conf.Security.RateLimit.Burst < 1 {
		errors = append(errors, "rate limit burst must be at least 1")
	}

	// Validate logging configuration
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[conf.Logging.Level] {
		errors = append(errors, "logging level must be one of: debug, info, warn, error")
	}

	validLogFormats := map[string]bool{
		"json": true, "text": true,
	}
	if !validLogFormats[conf.Logging.Format] {
		errors = append(errors, "logging format must be one of: json, text")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration errors:\n- %s", errors[0])
	}

	return nil
}
