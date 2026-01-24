package service

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Get all services
		api.GET("/services", GetServices)
		
		// Get services by part (1 or 2)
		api.GET("/services/part", GetServicesByPartHandler)
		
		// Create service
		api.POST("/services", CreateServiceHandler)
		
		// Update service
		api.PUT("/services/:id", UpdateServiceHandler)
		
		// Delete service
		api.DELETE("/services/:id", DeleteServiceHandler)
	}
}
