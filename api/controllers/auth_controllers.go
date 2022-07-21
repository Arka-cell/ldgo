package controllers

import (
	"net/http"

	"github.com/Arka-cell/ldgo/api/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) signUp(c *gin.Context) {
	var shop = models.Shop{}
	if err := c.BindJSON(&shop); err != nil {
		return
	}
	shop.Prepare()
	err := shop.Validate("create")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "inconsistent data from your request"})
		return
	}
	shopCreated, err := shop.SaveShop(server.DB)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "unable to create your shop"})
		return
	}
	if shopCreated != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": shopCreated.Title + " created sucessfully!"})
		return
	}

}

func login(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{})
}
