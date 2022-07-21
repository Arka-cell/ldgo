package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Server struct {
	DB *gorm.DB
}

func (s *Server) Router() *gin.Engine {
	return gin.Default()
}
