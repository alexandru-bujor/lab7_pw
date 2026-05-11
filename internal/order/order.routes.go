package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	handler := NewHandler()

	orders := router.Group("/api/orders")
	{
		orders.GET("", handler.List)
		orders.GET("/:id", handler.GetByID)
		orders.POST("", handler.Create)
		orders.PUT("/:id/status", handler.UpdateStatus)
		orders.POST("/:id/mark-sold", handler.MarkAsSold)
		orders.DELETE("/:id", handler.Delete)
		orders.GET("/client/:clientId", handler.GetByClientID)
		orders.GET("/secondary", handler.GetSecondaryOrders)
	}
}
