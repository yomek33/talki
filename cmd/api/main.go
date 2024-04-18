package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/yomek33/talki/internal/handler"
	"github.com/yomek33/talki/internal/repository"
	"github.com/yomek33/talki/internal/repository/dbrepo"
)

const port = 8080

type application struct {
	DB       *gorm.DB
	UserRepo *repository.UserRepository
	ArticleRepo *repository.ArticleRepository
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

	// JWT middlewareの設定
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte("your_secret_key"), 
	}))


	userRepo := dbrepo.NewUserRepo(app.DB)
	userHandler := handler.NewUserHandler(userRepo)
	articleRepo := dbrepo.NewArticleRepo(app.DB)
	articleHandler := handler.NewArticleHandler(articleRepo)


	// ルーティング
	userHandler.HandleUsers(e)
	articleHandler.HandleArticles(e)

	// サーバーを起動
	log.Println("Server is running at port", port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}