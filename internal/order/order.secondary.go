package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetSecondaryOrders(c *gin.Context) {
	orders, err := h.svc.GetSecondaryOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
