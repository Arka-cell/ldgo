package api

import (
	"fmt"
	"log"
	"os"

	"github.com/Arka-cell/ldgo/api/controllers"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}

func Run() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	println(os.Getenv("DB_USER") + " is now connected" + " to " + os.Getenv("DB_NAME"))
	server.Migrate(
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	router := server.InitializeRouter()
	router.Run("0.0.0.0:8080")
}
