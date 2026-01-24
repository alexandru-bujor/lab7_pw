package order

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

// List returns all orders
func (h *Handler) List(c *gin.Context) {
	status := c.Query("status")
	orders, err := h.svc.List(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("📊 Fetched %d orders", len(orders))
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// GetByID returns a specific order
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// Create creates a new order (from website or manual entry)
func (h *Handler) Create(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("❌ Order creation validation error: %v", err)
		log.Printf("📦 Received order data: %+v", input)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "details": "Validation failed. Please check all required fields are provided."})
		return
	}
	
	log.Printf("📦 Creating order: Client=%s, ProductID=%d, Quantity=%d, Price=%.2f", 
		input.ClientName, input.ProductID, input.Quantity, input.SellPrice)

	order, err := h.svc.CreateFromWebsite(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send Telegram notification
	telegram.SendOrderNotification(
		order.ID,
		order.ClientName,
		order.ClientPhone,
		order.ProductTitle,
		order.SellPrice,
	)

	c.JSON(http.StatusCreated, gin.H{"order": order})
}

// UpdateStatus updates the status of an order
func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.svc.UpdateStatus(id, body.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// MarkAsSold marks an order as sold and links to inventory item
func (h *Handler) MarkAsSold(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var body struct {
		InventoryItemID int     `json:"inventory_item_id" binding:"required"`
		SellPrice       *float64 `json:"sell_price"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.svc.MarkAsSold(id, body.InventoryItemID, body.SellPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// GetByClientID returns all orders for a specific client
func (h *Handler) GetByClientID(c *gin.Context) {
	clientID, err := strconv.Atoi(c.Param("clientId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	orders, err := h.svc.GetByClientID(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// Delete deletes an order
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted successfully"})
}
