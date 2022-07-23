package controllers

import (
	"net/http"

	"github.com/Arka-cell/ldgo/api/auth"
	"github.com/Arka-cell/ldgo/api/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) createProduct(c *gin.Context) {
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	var product = models.Product{}
	var shop = models.Shop{}

	if err := c.BindJSON(&product); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to bind data to product"})
		return
	}
	shop.FindShopByID(s.DB, uid)

	product.ShopID = uid
	product.Shop = shop
	product.Prepare()
	product.SaveProduct(s.DB)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Product Created"})

}
