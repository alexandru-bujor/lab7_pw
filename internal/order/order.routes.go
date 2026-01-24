package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	handler := NewHandler()
	
	orders := router.Group("/api/orders")
	{
		orders.GET("", handler.List)                      // Get all orders
		orders.GET("/:id", handler.GetByID)               // Get order by ID
		orders.POST("", handler.Create)                   // Create order from website/client
		orders.PUT("/:id/status", handler.UpdateStatus)   // Update order status
		orders.POST("/:id/mark-sold", handler.MarkAsSold) // Mark order as sold (link inventory item)
		orders.DELETE("/:id", handler.Delete)             // Delete order
		orders.GET("/client/:clientId", handler.GetByClientID) // Get orders for a client
		orders.GET("/secondary", handler.GetSecondaryOrders)   // Get secondary panel orders
	}
}
