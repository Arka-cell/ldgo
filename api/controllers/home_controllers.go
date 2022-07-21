package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getHome(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, gin.H{"message": "Welcome"})
}
