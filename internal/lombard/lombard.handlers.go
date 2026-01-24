package lombard

import (
	"log"
	"net/http"
	"strconv"

	"MegaMobileBack/internal/telegram"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

// Create creates a new lombard request
func (h *Handler) Create(c *gin.Context) {
	var input CreateLombardRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.svc.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("📝 Created lombard request #%d from %s", request.ID, request.CustomerName)
	
	// Send Telegram notification
	itemsCount := len(request.Products)
	telegram.SendLombardRequestNotification(
		request.ID,
		request.CustomerName,
		request.CustomerPhone,
		request.RequestType,
		itemsCount,
	)
	
	c.JSON(http.StatusCreated, gin.H{"request": request})
}

// List returns all lombard requests
func (h *Handler) List(c *gin.Context) {
	requests, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("📊 Fetched %d lombard requests", len(requests))
	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// GetByID returns a specific lombard request
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	request, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lombard request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request": request})
}

// UpdateStatus updates the status of a lombard request
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	var body UpdateStatusInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.svc.UpdateStatus(id, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Updated lombard request #%d status to %s", id, body.Status)
	c.JSON(http.StatusOK, gin.H{"request": request})
}

// Delete deletes a lombard request
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("🗑️ Deleted lombard request #%d", id)
	c.JSON(http.StatusOK, gin.H{"message": "lombard request deleted successfully"})
}

