package category

import "github.com/gin-gonic/gin"

// Routes registers category CRUD endpoints under /api/categories.
func Routes(r *gin.RouterGroup) {
	// list / create
	r.GET("/", GetCategoriesHandler)
	r.POST("/", CreateCategoryHandler)
	// Support both with-trailing-slash and without-trailing-slash paths
	r.GET("", GetCategoriesHandler)
	r.POST("", CreateCategoryHandler)

	// update / delete
	r.PUT("/:id", UpdateCategoryHandler)
	r.DELETE("/:id", DeleteCategoryHandler)

	// category template endpoints
	r.GET("/:id/template", GetCategoryTemplateHandler)
	r.PUT("/:id/template", UpsertCategoryTemplateHandler)
	r.DELETE("/:id/template", DeleteCategoryTemplateHandler)
}
