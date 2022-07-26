package controllers

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) InitializeRouter() *gin.Engine {
	router := s.Router()
	router.GET("/", s.getHome)
	router.POST("/api/v0/signup", s.signUp)
	router.POST("/api/v0/login", s.login)
	router.PATCH("/api/v0/shops", s.partialUpdateShop)
	router.DELETE("/api/v0/shops", s.deleteShop)
	router.GET("/api/v0/shops", s.getAllShops)
	router.GET("/api/v0/shops/:id", s.getShop)
	router.POST("/api/v0/products", s.createProduct)
	router.GET("/api/v0/products", s.getAllProducts)
	router.GET("/api/v0/products/:id", s.getProduct)
	router.PATCH("/api/v0/products/:id", s.partialUpdateProduct)
	router.DELETE("/api/v0/products/:id", s.deleteProduct)
	return router
}
