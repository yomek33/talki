package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yomek33/talki/internal/config"
	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/handler"
	"github.com/yomek33/talki/internal/services"
	"github.com/yomek33/talki/internal/stores"
)

type application struct {
	DB          *gorm.DB
	GeminiClient *gemini.Client
}

func main() {
	// Initialize Echo
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&tls=%s&parseTime=True&loc=Local",
		cfg.TiDBUser,
		cfg.TiDBPassword,
		cfg.TiDBHost,
		cfg.TiDBPort,
		cfg.TiDBDBName,
		cfg.UseSSL,
	)

	// Connect to the database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	// Initialize Gemini client
	geminiClient, err := gemini.NewClient(context.Background(), cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	defer geminiClient.Close()
	// Initialize application structure
	app := &application{DB: db, GeminiClient: geminiClient}
	// Initialize stores, services, and handlers
	stores := stores.NewStores(app.DB)
	services := services.NewServices(stores, app.GeminiClient)
	h := handler.NewHandler(services)

	// Set up routes
	h.SetDefault(e)
	h.SetAPIRoutes(e)

	// Start the server
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
