// Copyright 2026 Mano
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mano/mucal"
	"github.com/mano/mucal/internal/api"
	"github.com/mano/mucal/internal/config"
	"github.com/mano/mucal/internal/version"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	// If a positional argument is provided, use it as config path
	if len(flag.Args()) > 0 {
		*configPath = flag.Args()[0]
	}

	log.Printf("μCal version: %s", version.Version)

	// Load configuration
	log.Printf("Loading configuration from: %s", *configPath)
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Loaded %d calendar(s)", len(cfg.Calendars))

	// Create API handler
	handler, err := api.NewHandler(cfg)
	if err != nil {
		log.Fatalf("Failed to create API handler: %v", err)
	}

	// Setup routes
	mux := http.NewServeMux()
	handler.SetupRoutes(mux)

	// Serve embedded frontend
	webFS := mucal.GetWebFS()
	fileServer := http.FileServer(webFS)
	mux.Handle("/", fileServer)

	// Wrap with middleware
	wrappedHandler := api.RecoveryMiddleware(
		api.LoggingMiddleware(
			api.CORSMiddleware(mux),
		),
	)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:         addr,
		Handler:      wrappedHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped")
}
