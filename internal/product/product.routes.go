package product

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	r.GET("/", GetProducts)
	r.GET("", GetProducts)
	r.POST("/", CreateProduct)
	r.POST("", CreateProduct)
	r.GET("/secondary", GetSecondaryProducts)
	r.GET("/slug/:slug", GetProductBySlugHandler)

	r.GET("/:id/variants", GetProductVariantsHandler)
	r.PUT("/:id", UpdateProductHandler)
	r.DELETE("/:id", DeleteProductHandler)
}
