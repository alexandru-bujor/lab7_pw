package make

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/makes", CreateMakeHandler)
		api.GET("/makes", GetAllMakesHandler)
		api.PUT("/makes/:id", UpdateMakeHandler)
		api.DELETE("/makes/:id", DeleteMakeHandler)

		api.POST("/models", CreateModelHandler)
		api.PUT("/models/:id", UpdateModelHandler)
		api.DELETE("/models/:id", DeleteModelHandler)
	}
}
