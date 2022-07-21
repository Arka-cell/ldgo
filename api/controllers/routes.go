package controllers

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) InitializeRouter() *gin.Engine {
	router := s.Router()
	router.GET("/", getHome)
	return router
}
