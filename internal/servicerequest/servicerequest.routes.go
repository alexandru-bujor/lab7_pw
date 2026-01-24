package servicerequest

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	handler := NewHandler()

	serviceRequests := router.Group("/api/service-requests")
	{
		serviceRequests.GET("", handler.List)                    // Get all requests
		serviceRequests.GET("/:id", handler.GetByID)             // Get request by ID
		serviceRequests.POST("", handler.Create)                 // Create new request
		serviceRequests.PUT("/:id/status", handler.UpdateStatus) // Update request status
		serviceRequests.DELETE("/:id", handler.Delete)           // Delete request
	}
}

