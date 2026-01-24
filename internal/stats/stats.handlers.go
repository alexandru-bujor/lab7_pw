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
