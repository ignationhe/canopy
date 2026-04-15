// Package main is the entry point for the Canopy node.
// Canopy is a blockchain network implementation in Go.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/canopy-network/canopy/config"
	"github.com/canopy-network/canopy/node"
)

const (
	// Version is the current version of the Canopy node.
	Version = "0.1.0"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "config.json", "path to the configuration file")
	// Changed default dir to match XDG-style convention on my machine
	dataDir := flag.String("data-dir", ".canopy-data", "path to the data directory")
	logLevel := flag.String("log-level", "debug", "log level (debug, info, warn, error)")
	printVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	// Print version if requested
	if *printVersion {
		fmt.Printf("Canopy Node v%s\n", Version)
		os.Exit(0)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "[canopy] ", log.LstdFlags|log.Lshortfile)
	logger.Printf("Starting Canopy Node v%s", Version)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Printf("Warning: could not load config from %s: %v — using defaults", *configPath, err)
		cfg = config.Default()
	}

	// Apply flag overrides
	if *dataDir != ".canopy-data" {
		cfg.DataDir = *dataDir
	}
	if *logLevel != "debug" {
		cfg.LogLevel = *logLevel
	}

	// Ensure data directory exists
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		logger.Fatalf("Failed to create data directory %s: %v", cfg.DataDir, err)
	}

	// Create and start the node
	n, err := node.New(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to create node: %v", err)
	}

	if err := n.Start(); err != nil {
		logger.Fatalf("Failed to start node: %v", err)
	}

	logger.Println("Node started successfully")

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	logger.Printf("Received signal %s, shutting down...", sig)

	// Graceful shutdown
	if err := n.Stop(); err != nil {
		logger.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}

	logger.Println("Node stopped gracefully")
}
