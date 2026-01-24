package servicecategory

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	       api := router.Group("/api")
	       {
		       api.POST("/service-categories", CreateServiceCategoryHandler)
		       api.GET("/service-categories", GetAllServiceCategoriesHandler)
		       api.PUT("/service-categories/:id", UpdateServiceCategoryHandler)
		       api.DELETE("/service-categories/:id", DeleteServiceCategoryHandler)
	       }
}
