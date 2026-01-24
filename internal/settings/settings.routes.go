package settings

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/secondary-panel/verify-password", VerifySecondaryPasswordHandler)
		api.POST("/secondary-panel/change-password", ChangeSecondaryPasswordHandler)
	}
}
