package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	handler := NewHandler()

	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("", handler.List)
			users.GET("/:id", handler.GetByID)
			users.POST("", handler.Create)
			users.PUT("/:id", handler.Update)
			users.DELETE("/:id", handler.Delete)
		}
	}
}
