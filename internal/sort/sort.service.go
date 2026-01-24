package sort

import "github.com/gin-gonic/gin"

func Sort(r *gin.RouterGroup) {

	r.GET("/", GetProducts)

}
