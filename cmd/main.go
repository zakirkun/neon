package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zakirkun/neon/internal/cli"
	"github.com/zakirkun/neon/internal/config"
	"github.com/zakirkun/neon/internal/docker"
	"github.com/zakirkun/neon/internal/docker/health"
	"github.com/zakirkun/neon/internal/logger"
)

var (
	version = "0.1.0"
	commit  = "development"
)

func main() {
	// Add version flag
	showVersion := flag.Bool("version", false, "Show version information")
	// Parse all flags
	flag.Parse()

	if *showVersion {
		fmt.Printf("Neon v%s (%s)\n", version, commit)
		return
	}

	// Initialize logger
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	logDir := filepath.Join(homeDir, ".neon", "logs")
	if err := logger.Init(logDir); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize config
	if err := config.Load(); err != nil {
		if !os.IsNotExist(err) {
			logger.Error(err, "Failed to load configuration")
		} else {
			logger.Warn("No config file found, using defaults")
		}
	}

	// Initialize Docker client
	client, err := docker.NewClient()
	if err != nil {
		logger.Error(err, "Failed to initialize Docker client")
		os.Exit(1)
	}

	// Check Docker and Swarm status
	checker := health.NewChecker(client.Client)

	if err := checker.CheckDockerStatus(); err != nil {
		logger.Error(err, "Docker health check failed")
		fmt.Println("Error: Docker daemon is not running")
		os.Exit(1)
	}

	if err := checker.CheckSwarmStatus(); err != nil {
		logger.Error(err, "Swarm health check failed")
		fmt.Println("Warning: Docker Swarm mode is not enabled")
		fmt.Println("Some features may not be available")
		fmt.Println("To enable Swarm mode, run: docker swarm init")
	}

	if err := cli.Execute(); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
