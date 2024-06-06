package application

import (
  "os"
  "strconv"
)

// Config struct holds the configuration settings for the application.
type Config struct {
	RedisAddress string
	ServerPort   uint16
}

// LoadConfig loads the configuration from environment variables, falling back to default values if they are not set.
func LoadConfig() Config {
	// Initialize the Config struct with default values
	cfg := Config{
		RedisAddress: "localhost:6379", // Default Redis address
		ServerPort:   3000,				// Default server port
	}

	// Check if the REDIS_ADDR environment variable is set
	if redisAddr, exists := os.LookupEnv("REDIS_ADDR"); exists {
		cfg.RedisAddress = redisAddr // Update Redis address from the environment variable
	}

	// Check if the SERVER_PORT environment variable is set
	if serverPort, exists := os.LookupEnv("SERVER_PORT"); exists {
		// Parse the server port from the environment variable as an unsigned integer
		if port, err := strconv.ParseUint(serverPort, 10, 16); err == nil {
		cfg.ServerPort = uint16(port) // Update server port from the parsed value
		}
	}

	// Return the final configuration
	return cfg
}