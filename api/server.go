package api

import "github.com/Arka-cell/ldgo/api/controllers"

var server = controllers.Server{}

func Run() {
	router := server.Router()
	router.Run("localhost:8080")
}
