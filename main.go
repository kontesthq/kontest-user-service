package main

import (
	"github.com/kontesthq/go-consul-service-manager/consulservicemanager"
	"io"
	"kontest-user-service/database"
	routes "kontest-user-service/route"
	"kontest-user-service/utils"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func initializeVariables() {
	// Get the hostname of the machine
	hostname, err := os.Hostname()
	if err != nil {
		slog.Error("Error fetching hostname", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Attempt to read the KONTEST_API_SERVER_HOST environment variable
	if host := os.Getenv("KONTEST_USER_SERVICE_HOST"); host != "" {
		utils.ApplicationHost = host // Override with the environment variable if set
	} else {
		utils.ApplicationHost = hostname // Use the machine's hostname if the env var is not set
	}

	// Attempt to read the KONTEST_API_SERVER_PORT environment variable
	if port := os.Getenv("KONTEST_USER_SERVICE_PORT"); port != "" {
		parsedPort, err := strconv.Atoi(port)
		if err != nil {
			slog.Error("Invalid port value", slog.String("error", err.Error()), slog.String("port", port))
			os.Exit(1) // Exit the program with a non-zero status code
		}
		utils.ApplicationPort = parsedPort // Override with the environment variable if set
		slog.Info("Application port set from environment variable", slog.Int("applicationPort", utils.ApplicationPort))
	}

	// Attempt to read the CONSUL_ADDRESS environment variable
	if host := os.Getenv("CONSUL_HOST"); host != "" {
		utils.ConsulHost = host // Override with the environment variable if set
	}

	// Attempt to read the CONSUL_PORT environment variable
	if port := os.Getenv("CONSUL_PORT"); port != "" {
		if portInt, err := strconv.Atoi(port); err == nil {
			utils.ConsulPort = portInt // Override with the environment variable if set and valid
		}
	}

	// Attempt to read the DB_HOST environment variable
	if host := os.Getenv("DB_HOST"); host != "" {
		utils.DbHost = host // Override with the environment variable if set
	}

	// Attempt to read the DB_PORT environment variable
	if port := os.Getenv("DB_PORT"); port != "" {
		utils.DbPort = port // Override with the environment variable if set
	}

	// Attempt to read the DB_NAME environment variable
	if name := os.Getenv("DB_NAME"); name != "" {
		utils.DbName = name // Override with the environment variable if set
	}

	// Attempt to read the DB_USER environment variable
	if user := os.Getenv("DB_USER"); user != "" {
		utils.DbUser = user // Override with the environment variable if set
	}

	// Attempt to read the DB_PASSWORD environment variable
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		utils.DbPassword = password // Override with the environment variable if set
	}

	// Attempt to read the DB_SSL_MODE environment variable
	if sslMode := os.Getenv("DB_SSL_MODE"); sslMode != "" {
		utils.IsSSLModeEnabled = sslMode == "enable"
	}
}

func setupLogging() *os.File {
	LOG_FILE := os.Getenv("LOG_FILE")

	if LOG_FILE == "" {
		LOG_FILE = "tmp/logs/logs.log"
	}

	// Get the directory from the log file path
	logDir := filepath.Dir(LOG_FILE)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		slog.Error("Failed to create log directory", slog.String("error", err.Error()))
		os.Exit(1)
	}

	handlerOptions := &slog.HandlerOptions{
		AddSource: true,
	}
	// Open or create a log file
	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// Handle error if the log file cannot be opened or created
		slog.Error("Failed to open log file", slog.String("error", err.Error()))
		os.Exit(1)
	}

	w := io.MultiWriter(os.Stdout, logFile)

	// Configure slog to output JSON
	slog.SetDefault(slog.New(slog.NewJSONHandler(w, handlerOptions)))

	// Return the log file to close it in the main function
	return logFile
}

func main() {
	initializeVariables()

	logFile := setupLogging()
	// Ensure the log file is closed when the program exits
	defer logFile.Close()

	// Log server restart with a timestamp
	slog.Info("Server restarted", slog.Time("time", time.Now()))

	// Initialize the database connection
	database.InitializeDatabase(utils.DbName, utils.DbPort, utils.DbHost, utils.DbUser, utils.DbPassword, map[bool]string{true: "enable", false: "disable"}[utils.IsSSLModeEnabled])
	database.SetupDatabase()
	defer database.CloseDB()

	consulService := consulservicemanager.NewConsulService(utils.ConsulHost, utils.ConsulPort)
	consulService.Start(utils.ApplicationHost, utils.ApplicationPort, utils.ServiceName, []string{})

	router := http.NewServeMux()

	routes.RegisterRoutes(router)

	server := http.Server{
		Addr:    ":" + strconv.Itoa(utils.ApplicationPort),
		Handler: router,
	}

	slog.Info("Server listening", slog.Int("port", utils.ApplicationPort))

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
