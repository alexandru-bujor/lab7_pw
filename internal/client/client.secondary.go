package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSecondaryClients returns only clients visible in secondary panel
func (h *Handler) GetSecondaryClients(c *gin.Context) {
	clients, err := h.service.GetSecondaryClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, clients)
}
