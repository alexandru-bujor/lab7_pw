package client

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	handler := NewHandler()

	clients := router.Group("/api/clients")
	{
		clients.GET("", handler.GetAllClients)
		clients.GET("/:id", handler.GetClientByID)
		clients.POST("", handler.CreateClient)
		clients.PUT("/:id", handler.UpdateClient)
		clients.DELETE("/:id", handler.DeleteClient)
		clients.POST("/:id/orders", handler.RecordOrder)
		clients.POST("/get-or-create", handler.GetOrCreateClientByEmail)
		clients.GET("/secondary", handler.GetSecondaryClients) // Get secondary panel clients
	}
}
