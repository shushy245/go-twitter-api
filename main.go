// main.go
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"hello/tweet"
	"hello/user"
	"os"
)

var db *gorm.DB
var err error

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBName     string
	DBPassword string
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}
}

func main() {
	cfg := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBName:     os.Getenv("DB_NAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
	}

	db, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword))

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connected to the database")
	}

	defer db.Close()

	db.AutoMigrate(&user.User{}, &tweet.Tweet{})

	router := gin.Default()

	user.Setup(router, db)
	tweet.Setup(router, db)

	router.Run(":9000")
}
