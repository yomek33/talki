package main

import (
	"fmt"
	"log"

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

	dsn := "root:password@tcp(mysql:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"
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