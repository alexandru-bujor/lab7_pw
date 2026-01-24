package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	handler := NewHandler()

	api := router.Group("/api")
	{
		// User management routes
		users := api.Group("/users")
		{
			users.GET("", handler.List)       // Get all users
			users.GET("/:id", handler.GetByID) // Get user by ID
			users.POST("", handler.Create)     // Create new user
			users.PUT("/:id", handler.Update)  // Update user
			users.DELETE("/:id", handler.Delete) // Delete user
		}
	}
}

