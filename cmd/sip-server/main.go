package main

import (
    "flag"
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/e173-gateway/e173_go_gateway/pkg/sip"
)

func main() {
    port := flag.Int("port", 5060, "SIP server port")
    whatsappKey := flag.String("whatsapp-key", "", "WhatsApp Business API key")
    flag.Parse()

    if *whatsappKey == "" {
        log.Println("Warning: No WhatsApp API key provided, validation will be limited")
    }

    log.Printf("Starting E173 SIP Gateway Server on port %d", *port)
    
    // Create SIP server
    server := sip.NewBasicSIPServer(*port, *whatsappKey)
    
    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Start server in goroutine
    go func() {
        if err := server.Start(); err != nil {
            log.Fatalf("Failed to start SIP server: %v", err)
        }
    }()
    
    // Wait for shutdown signal
    <-sigChan
    log.Println("Shutting down SIP server...")
}
