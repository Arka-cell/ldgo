package controllers

import (
	"net/http"
	"strconv"

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

func (s *Server) getAllProducts(c *gin.Context) {
	product := models.Product{}
	products, err := product.FindAllProducts(s.DB)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"products": products})
}

func (s *Server) getProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	uid := uint64(id)
	if err != nil {
		return
	}
	product := models.Product{}
	product.FindProductByID(s.DB, uid)
	c.IndentedJSON(http.StatusOK, gin.H{"product": product})
}

func (s *Server) partialUpdateProduct(c *gin.Context) {
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	pk, err := strconv.Atoi(c.Param("id"))
	pk32 := uint64(pk)
	if err != nil {
		return
	}
	product := models.Product{}
	product.FindProductByID(s.DB, pk32)
	if product.ShopID != uid {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	if err := c.BindJSON(&product); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to bind JSON to credentials struct"})
		return
	}
	pk64 := int32(pk32)
	updatedProduct, err := product.PartialUpdateProduct(s.DB, pk64, product.Name, product.Description, product.Price)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"product": updatedProduct})
}
