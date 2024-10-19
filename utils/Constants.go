package utils

import (
	"os"
	"strconv"
)

const (
	UserIdRequestHeader = "logged-in-user-id"
)

var (
	API_GATEWAY_HOST = "localhost"
	API_GATEWAY_PORT = 5153
)

func InitializeVariables() {
	// Get the API Gateway host and port from the environment variables
	if host := os.Getenv("API_GATEWAY_HOST"); host != "" {
		API_GATEWAY_HOST = host
	}

	if port := os.Getenv("API_GATEWAY_PORT"); port != "" {
		API_GATEWAY_PORT, _ = strconv.Atoi(port)
	}
}
