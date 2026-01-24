package inventory

import "github.com/gin-gonic/gin"

func RegisterRoutes(api *gin.RouterGroup) {
	h := NewHandler()
	inv := api.Group("/inventory")
	{
		inv.GET("", h.List)
		inv.POST("", h.Create)
		inv.PUT("/:id", h.Update)
		inv.DELETE("/:id", h.Delete)
		inv.POST("/:id/assign-client", h.AssignClient)
		inv.POST("/:id/mark-sold", h.MarkSold)
		inv.POST("/sync", h.SyncToInventory)		// Sections management
		inv.GET("/sections", h.ListSections)
		inv.POST("/sections", h.CreateSection)
		inv.PUT("/sections/:id", h.UpdateSection)
		inv.DELETE("/sections/:id", h.DeleteSection)	}
}
