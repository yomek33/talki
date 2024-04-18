package main

package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

)

const port = 8080

type application struct {
	DB       *gorm.DB
	TaskRepo repository.TaskRepository
}

func main() {
	var app application
	var err error

	dsn := "root:password@tcp(mysql:3306)/todo?charset=utf8mb4&parseTime=True&loc=Local"

	// GORM DB接続
	app.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect DB: ", err)
	}

	// httpサーバーの起動
	e := app.routes()

	// サーバーの開始
	log.Println("Starting server on port", port)
	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}