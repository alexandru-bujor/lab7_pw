package product

import (
	"MegaMobileBack/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSecondaryProducts(c *gin.Context) {
	var products []Product
	err := db.DB.
		Preload("Images").
		Find(&products).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
