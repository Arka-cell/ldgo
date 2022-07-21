package api

import "github.com/Arka-cell/ldgo/api/controllers"

var server = controllers.Server{}

func Run() {
	router := server.InitializeRouter()
	router.Run("localhost:8080")
}
