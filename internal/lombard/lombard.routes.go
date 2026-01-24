package lombard

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	handler := NewHandler()

	lombard := router.Group("/api/lombard")
	{
		requests := lombard.Group("/requests")
		{
			requests.GET("", handler.List)                    // Get all requests
			requests.GET("/:id", handler.GetByID)             // Get request by ID
			requests.POST("", handler.Create)                 // Create new request
			requests.PUT("/:id/status", handler.UpdateStatus) // Update request status
			requests.DELETE("/:id", handler.Delete)           // Delete request
		}
	}
}

