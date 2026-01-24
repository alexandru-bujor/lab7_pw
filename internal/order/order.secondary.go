package order

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSecondaryOrders returns only orders visible in secondary panel
func (h *Handler) GetSecondaryOrders(c *gin.Context) {
	orders, err := h.svc.GetSecondaryOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("📊 Fetched %d secondary orders", len(orders))
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
