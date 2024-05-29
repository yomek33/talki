package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/handler"
	"github.com/yomek33/talki/services"
	"github.com/yomek33/talki/stores"
)

const port = 8080

type application struct {
	DB       *gorm.DB
}


func main() {
	e := echo.New()
	var app application
	var err error

    tidbUser := os.Getenv("TIDB_USER")
    tidbPassword := os.Getenv("TIDB_PASSWORD")
    tidbHost := os.Getenv("TIDB_HOST")
    tidbPort := os.Getenv("TIDB_PORT")
    tidbDBName := os.Getenv("TIDB_DB_NAME")
    useSSL := os.Getenv("USE_SSL")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&tls=%s&parseTime=True&loc=Local", tidbUser, tidbPassword, tidbHost, tidbPort, tidbDBName, useSSL)


	// GORM DB接続
	app.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB: ", err)
	}

	stores := stores.NewStores(app.DB)

	services := services.NewServices(stores)

	h := handler.NewHandler(services)

	h.SetDefault(e)

	h.SetAPIRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}