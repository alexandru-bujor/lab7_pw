package stats

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
    handler := NewHandler()
    api.GET("/stats", handler.Summary)
}
