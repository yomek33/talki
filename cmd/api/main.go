package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yomek33/talki/internal/config"
	"github.com/yomek33/talki/internal/gemini"
	"github.com/yomek33/talki/internal/handler"
	"github.com/yomek33/talki/internal/services"
	"github.com/yomek33/talki/internal/stores"
)

type application struct {
	DB           *gorm.DB
	GeminiClient *gemini.Client
	Firebase     *handler.Firebase
}

func main() {
	// Initialize Echo
	e := handler.Echo()

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

	geminiClient, err := gemini.NewClient(context.Background(), cfg.GeminiAPIKey)
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	defer geminiClient.Close()

	firebaseInstance, err := handler.InitFirebase(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize application structure
	app := &application{
		DB:           db,
		GeminiClient: geminiClient,
		Firebase:     firebaseInstance,
	}

	stores := stores.NewStores(app.DB)
	services := services.NewServices(stores, app.GeminiClient)
	h := handler.NewHandler(services, cfg.JWTSecretKey, app.Firebase)

	e.Use(handler.FirebaseAuthMiddleware(app.Firebase.AuthClient))
	h.SetDefault(e)
	h.SetAPIRoutes(e)

	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
