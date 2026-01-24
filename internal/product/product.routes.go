package product

import "github.com/gin-gonic/gin"

func Routes(r *gin.RouterGroup) {
	// Support both with-trailing-slash and without-trailing-slash paths
	r.GET("/", GetProducts)
	r.GET("", GetProducts)

	r.POST("/", CreateProduct)
	r.POST("", CreateProduct)

	r.GET("/:id/variants", GetProductVariantsHandler)
	r.PUT("/:id", UpdateProductHandler)
	r.DELETE("/:id", DeleteProductHandler)
	
	// Secondary panel endpoint
	r.GET("/secondary", GetSecondaryProducts)
}
