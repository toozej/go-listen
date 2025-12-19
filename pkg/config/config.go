// Package config provides secure configuration management for the go-listen application.
//
// This package handles loading configuration from environment variables and .env files
// with built-in security measures to prevent path traversal attacks. It uses the
// github.com/caarlos0/env library for environment variable parsing and
// github.com/joho/godotenv for .env file loading.
//
// The configuration loading follows a priority order:
//  1. Environment variables (highest priority)
//  2. .env file in current working directory
//  3. Default values (if any)
//
// Security features:
//   - Path traversal protection for .env file loading
//   - Secure file path resolution using filepath.Abs and filepath.Rel
//   - Validation against directory traversal attempts
//
// Example usage:
//
//	import "github.com/toozej/go-listen/pkg/config"
//
//	func main() {
//		conf := config.GetEnvVars(false)
//		fmt.Printf("Server: %s\n", conf.Server.Address())
//	}
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Config represents the application configuration structure.
//
// This struct defines all configurable parameters for the go-listen
// application. Fields are tagged with struct tags that correspond to
// environment variable names for automatic parsing.
//
// Configuration sections:
//   - Server: HTTP server configuration (host, port)
//   - Spotify: Spotify API credentials and settings
//   - Security: Security-related settings (rate limiting)
//   - Logging: Logging configuration (level, format, output)
//
// Example:
//
//	conf := config.GetEnvVars(false)
//	fmt.Printf("Server will run on: %s\n", conf.Server.Address())
type Config struct {
	Server   ServerConfig   `envPrefix:"SERVER_"`
	Spotify  SpotifyConfig  `envPrefix:"SPOTIFY_"`
	Security SecurityConfig `envPrefix:"SECURITY_"`
	Logging  LoggingConfig  `envPrefix:"LOGGING_"`
	Scraper  ScraperConfig  `envPrefix:"SCRAPER_"`
}

type ServerConfig struct {
	Host         string `env:"HOST" envDefault:"127.0.0.1"`
	Port         int    `env:"PORT" envDefault:"8080"`
	ReadTimeout  int    `env:"READ_TIMEOUT_SECONDS" envDefault:"30"`
	WriteTimeout int    `env:"WRITE_TIMEOUT_SECONDS" envDefault:"60"`
	IdleTimeout  int    `env:"IDLE_TIMEOUT_SECONDS" envDefault:"120"`
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

type ScraperConfig struct {
	TimeoutSeconds int    `env:"TIMEOUT_SECONDS" envDefault:"30"`
	MaxRetries     int    `env:"MAX_RETRIES" envDefault:"3"`
	RetryBackoff   int    `env:"RETRY_BACKOFF_SECONDS" envDefault:"2"`
	UserAgent      string `env:"USER_AGENT" envDefault:"go-listen/1.0 (Web Scraper)"`
	MaxContentSize int64  `env:"MAX_CONTENT_SIZE" envDefault:"10485760"` // 10MB in bytes
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

// GetEnvVars loads and returns the application configuration from environment
// variables and .env files with comprehensive security validation.
//
// This function performs the following operations:
//  1. Securely determines the current working directory
//  2. Constructs and validates the .env file path to prevent traversal attacks
//  3. Loads .env file if it exists in the current directory
//  4. Parses environment variables into the Config struct
//  5. Validates the configuration for correctness
//  6. Returns the populated configuration
//
// Security measures implemented:
//   - Path traversal detection and prevention using filepath.Rel
//   - Absolute path resolution for secure path operations
//   - Validation against ".." sequences in relative paths
//   - Safe file existence checking before loading
//
// The function will terminate the program with os.Exit(1) if any critical
// errors occur during configuration loading, such as:
//   - Current directory access failures
//   - Path traversal attempts detected
//   - .env file parsing errors
//   - Environment variable parsing failures
//   - Configuration validation failures
//
// Parameters:
//   - debug: If true, prints detailed configuration information
//
// Returns:
//   - Config: A populated and validated configuration struct
//
// Example:
//
//	// Load configuration with debug output
//	conf := config.GetEnvVars(true)
//
//	// Use configuration
//	server := &http.Server{Addr: conf.Server.Address()}
func GetEnvVars(debug bool) Config {
	// Get current working directory for secure file operations
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current working directory: %s\n", err)
		os.Exit(1)
	}

	// Construct secure path for .env file within current directory
	envPath := filepath.Join(cwd, ".env")

	// Ensure the path is within our expected directory (prevent traversal)
	cleanEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		fmt.Printf("Error resolving .env file path: %s\n", err)
		os.Exit(1)
	}
	cleanCwd, err := filepath.Abs(cwd)
	if err != nil {
		fmt.Printf("Error resolving current directory: %s\n", err)
		os.Exit(1)
	}
	relPath, err := filepath.Rel(cleanCwd, cleanEnvPath)
	if err != nil || strings.Contains(relPath, "..") {
		fmt.Printf("Error: .env file path traversal detected\n")
		os.Exit(1)
	}

	// Load .env file if it exists (will not override existing environment variables)
	if _, err := os.Stat(envPath); err == nil {
		if err := godotenv.Load(envPath); err != nil {
			fmt.Printf("Error loading .env file: %s\n", err)
			os.Exit(1)
		} else if debug {
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
	if conf.Server.ReadTimeout < 1 {
		errors = append(errors, "server read timeout must be at least 1 second")
	}
	if conf.Server.WriteTimeout < 1 {
		errors = append(errors, "server write timeout must be at least 1 second")
	}
	if conf.Server.IdleTimeout < 1 {
		errors = append(errors, "server idle timeout must be at least 1 second")
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

	// Validate scraper configuration
	if conf.Scraper.TimeoutSeconds < 1 {
		errors = append(errors, "scraper timeout must be at least 1 second")
	}
	if conf.Scraper.MaxRetries < 0 {
		errors = append(errors, "scraper max retries must be non-negative")
	}
	if conf.Scraper.RetryBackoff < 0 {
		errors = append(errors, "scraper retry backoff must be non-negative")
	}
	if conf.Scraper.MaxContentSize < 1 {
		errors = append(errors, "scraper max content size must be at least 1 byte")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration errors:\n- %s", errors[0])
	}

	return nil
}
