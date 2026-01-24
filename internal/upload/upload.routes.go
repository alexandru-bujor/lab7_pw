package upload

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup) {
	upload := rg.Group("/upload")
	{
		upload.POST("/product-image", UploadProductImageHandler)
	}
}
