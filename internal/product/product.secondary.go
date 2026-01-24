package product

import (
	"MegaMobileBack/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSecondaryProducts returns only products visible in secondary panel
func GetSecondaryProducts(c *gin.Context) {
	var products []Product
	err := db.DB.
		Where("is_visible_in_secondary = ?", true).
		Preload("Images").
		Preload("Variants").
		Find(&products).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
