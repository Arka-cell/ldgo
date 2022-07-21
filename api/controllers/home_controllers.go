package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getHome(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Welcome"})
}
