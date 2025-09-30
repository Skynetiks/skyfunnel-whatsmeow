package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"whatsmeow-service/config"
	"whatsmeow-service/handlers"
	"whatsmeow-service/services"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Database connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize services
	whatsAppService := services.NewWhatsAppMeowService(cfg, db)

	// Initialize handlers
	handlers := handlers.NewHandlers(cfg, whatsAppService)

	// Setup HTTP routes
	http.HandleFunc("/api/whatsmeow/send", handlers.SendMessage)
	http.HandleFunc("/api/whatsmeow/status", handlers.GetStatus)
	http.HandleFunc("/api/whatsmeow/qr", handlers.GetQR)
	http.HandleFunc("/api/whatsmeow/connect", handlers.Connect)
	http.HandleFunc("/api/whatsmeow/disconnect", handlers.Disconnect)
	http.HandleFunc("/health", handlers.Health)

	log.Printf("WhatsApp Meow service starting on port %d", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}

