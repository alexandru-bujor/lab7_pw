package category

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	r.GET("/", GetCategoriesHandler)
	r.POST("/", CreateCategoryHandler)
	r.GET("", GetCategoriesHandler)
	r.POST("", CreateCategoryHandler)

	r.PUT("/:id", UpdateCategoryHandler)
	r.DELETE("/:id", DeleteCategoryHandler)

	r.GET("/:id/template", GetCategoryTemplateHandler)
	r.PUT("/:id/template", UpsertCategoryTemplateHandler)
	r.DELETE("/:id/template", DeleteCategoryTemplateHandler)
}
