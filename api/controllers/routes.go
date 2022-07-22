package controllers

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) InitializeRouter() *gin.Engine {
	router := s.Router()
	router.GET("/", s.getHome)
	router.POST("/api/v0/signup", s.signUp)
	router.POST("/api/v0/login", s.login)
	router.PATCH("/api/v0/update-shop", s.partialUpdateShop)
	return router
}
