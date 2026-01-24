package inventory

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *Service }

func NewHandler() *Handler { return &Handler{svc: NewService()} }

func (h *Handler) List(c *gin.Context) {
	section := c.Query("section")
	status := c.Query("status")
	items, err := h.svc.List(section, status)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) Create(c *gin.Context) {
	var input InventoryItem
	if err := c.ShouldBindJSON(&input); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	if input.Section == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "section is required"}); return }
	if err := h.svc.Create(&input); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"item": input})
}

func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	var payload map[string]any
	if err := c.ShouldBindJSON(&payload); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := h.svc.Update(id, payload)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	if err := h.svc.Delete(id); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *Handler) AssignClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	var body struct{ ClientID int `json:"client_id" binding:"required"`}
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := h.svc.AssignClient(id, body.ClientID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *Handler) MarkSold(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	var body struct{ ClientID *int `json:"client_id"`; SellPrice *float64 `json:"sell_price"` }
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	item, err := h.svc.MarkSold(id, body.ClientID, body.SellPrice)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"item": item})
}

func (h *Handler) SyncToInventory(c *gin.Context) {
	var body struct{ Section string `json:"section"` }
	if err := c.ShouldBindJSON(&body); err != nil { body.Section = "Apple SH" }
	count, err := SyncAllProductsToInventory(body.Section)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"created": count, "message": "Inventory synced"})
}

func (h *Handler) ListSections(c *gin.Context) {
	sections, err := h.svc.ListSections()
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"sections": sections})
}

func (h *Handler) CreateSection(c *gin.Context) {
	var body struct{ Name string `json:"name" binding:"required"` }
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	sec, err := h.svc.CreateSection(body.Name)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"section": sec})
}

func (h *Handler) UpdateSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	var body struct{ Name string `json:"name" binding:"required"` }
	if err := c.ShouldBindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	sec, err := h.svc.UpdateSection(id, body.Name)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"section": sec})
}

func (h *Handler) DeleteSection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
	if err := h.svc.DeleteSection(id); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
