package main

import (
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/joho/godotenv"
    "github.com/e173-gateway/e173_go_gateway/pkg/sip"
    "github.com/e173-gateway/e173_go_gateway/pkg/database"
    "github.com/e173-gateway/e173_go_gateway/pkg/config"
)

func main() {
    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    
    // Load configuration
    cfg := config.LoadConfig()
    
    // Command line flags
    port := flag.Int("port", 5060, "SIP server port")
    whatsappKey := flag.String("whatsapp-key", "e42f7c9b-2a8e-4b86-a7e4-8f1de2c01f53", "WhatsApp API key")
    flag.Parse()

    if *whatsappKey == "" {
        log.Println("Warning: No WhatsApp API key provided, validation will be limited")
    }

    log.Printf("Starting E173 SIP Gateway Server on port %d", *port)
    
    // Initialize database connection
    dbPool, err := database.NewDBPool(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer dbPool.Close()
    log.Println("Successfully connected to database")
    
    // Create SIP server with database support
    server := sip.NewBasicSIPServerWithDB(*port, *whatsappKey, dbPool)
    
    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Start server in goroutine
    go func() {
        if err := server.Start(); err != nil {
            log.Fatalf("Failed to start SIP server: %v", err)
        }
    }()
    
    log.Printf("SIP server started successfully on port %d", *port)
    log.Println("WhatsApp API integration enabled with database caching")
    log.Println("Press Ctrl+C to shutdown...")
    
    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutting down SIP server...")
}