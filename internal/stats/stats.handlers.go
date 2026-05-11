package stats

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{ svc *Service }

func NewHandler() *Handler { return &Handler{svc: NewService()} }

func (h *Handler) Summary(c *gin.Context) {
	summary, err := h.svc.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": summary})
}

func (h *Handler) DailyStats(c *gin.Context) {
	today, yesterday, err := h.svc.GetDailyStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"today": today, "yesterday": yesterday})
}

func (h *Handler) History(c *gin.Context) {
	// For example, 30 days
	stats, err := h.svc.GetHistoricalDailyStats(30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"history": stats})
}
