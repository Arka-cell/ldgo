package controllers

import (
	"fmt"
	"net/http"

	"github.com/Arka-cell/ldgo/api/auth"
	"github.com/Arka-cell/ldgo/api/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

func (s *Server) login(c *gin.Context) {
	creds := credentials{}

	shop := models.Shop{}
	if err := c.BindJSON(&creds); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Authentication Failed"})
		fmt.Println("Unable to bind JSON to credentials struct")
		return
	}
	var err error

	err = s.DB.Debug().Model(models.Shop{}).Where("email = ?", creds.Email).Take(&shop).Error
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Authentication Failed"})
		return
	}
	err = models.VerifyPassword(shop.Password, creds.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Password or Email are false"})
		return
	}
	token, err := auth.CreateToken(shop.ID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}

func (s *Server) partialUpdateShop(c *gin.Context) {
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	shop := models.Shop{}
	shop.FindShopByID(s.DB, uid)
	if err != nil {
		return
	}
	if err := c.BindJSON(&shop); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Unable to bind JSON to credentials struct"})
		return
	}

	updatedShop, err := shop.PartialUpdateShop(s.DB, uid, shop.Password, shop.Title, shop.Email)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"shop": updatedShop})
}

func (s *Server) deleteShop(c *gin.Context) {
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	shop := models.Shop{}
	shop.FindShopByID(s.DB, uid)
	shop.DeleteShop(s.DB, uid)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Shop deleted"})
}
