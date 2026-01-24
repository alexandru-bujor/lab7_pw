package servicerequest

import (
	"fmt"
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

// Create creates a new service request
func (h *Handler) Create(c *gin.Context) {
	var input CreateServiceRequestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	request, err := h.svc.Create(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("📝 Created service request #%d from %s", request.ID, request.CustomerName)
	
	// Send Telegram notification
	deviceInfo := ""
	if request.DeviceMake != "" && request.DeviceModel != "" {
		deviceInfo = fmt.Sprintf("%s %s", request.DeviceMake, request.DeviceModel)
	} else if request.DeviceMake != "" {
		deviceInfo = request.DeviceMake
	} else if request.DeviceModel != "" {
		deviceInfo = request.DeviceModel
	} else {
		deviceInfo = "N/A"
	}
	
	telegram.SendServiceRequestNotification(
		request.ID,
		request.CustomerName,
		request.CustomerPhone,
		deviceInfo,
		request.ProblemDescription,
	)
	
	c.JSON(http.StatusCreated, gin.H{"request": request})
}

// List returns all service requests
func (h *Handler) List(c *gin.Context) {
	requests, err := h.svc.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("📊 Fetched %d service requests", len(requests))
	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// GetByID returns a specific service request
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request ID"})
		return
	}

	request, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "service request not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"request": request})
}

// UpdateStatus updates the status of a service request
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

	log.Printf("✅ Updated service request #%d status to %s", id, body.Status)
	c.JSON(http.StatusOK, gin.H{"request": request})
}

// Delete deletes a service request
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

	log.Printf("🗑️ Deleted service request #%d", id)
	c.JSON(http.StatusOK, gin.H{"message": "service request deleted successfully"})
}

