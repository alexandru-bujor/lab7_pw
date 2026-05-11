package router

import (
	"MegaMobileBack/internal/banners"
	"MegaMobileBack/internal/category"
	"MegaMobileBack/internal/client"
	"MegaMobileBack/internal/inventory"
	"MegaMobileBack/internal/lombard"
	"MegaMobileBack/internal/make"
	"MegaMobileBack/internal/order"
	"MegaMobileBack/internal/product"
	"MegaMobileBack/internal/service"
	servicecategory "MegaMobileBack/internal/servicecategory"
	"MegaMobileBack/internal/servicerequest"
	"MegaMobileBack/internal/settings"
	"MegaMobileBack/internal/stats"
	"MegaMobileBack/internal/upload"
	"MegaMobileBack/internal/user"
	"MegaMobileBack/pkg/db"

	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.RedirectTrailingSlash = false
	allowedOrigins := []string{
		"http://localhost:8081",
		"http://localhost:5173",
		"http://localhost:8080",
		"http://localhost:3000",
		"https://megamobile.md",
		"https://www.megamobile.md",
	}

	if additionalOrigins := os.Getenv("CORS_ORIGINS"); additionalOrigins != "" {
		for _, origin := range strings.Split(additionalOrigins, ",") {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello": "World",
		})
	})

	r.GET("/api/health", func(c *gin.Context) {
		dbOK := db.HealthCheck()

		status := http.StatusOK
		if !dbOK {
			status = http.StatusServiceUnavailable
		}

		c.JSON(status, gin.H{
			"status": "ok",
			"db": gin.H{
				"connected": dbOK,
			},
		})
	})

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.POST("/login", user.LoginHandler)
		productsGroup := api.Group("/products")
		product.Routes(productsGroup)
		categoriesGroup := api.Group("/categories")
		category.Routes(categoriesGroup)
		service.RegisterRoutes(r)
		servicecategory.RegisterRoutes(r)
		make.RegisterRoutes(r)
		upload.RegisterRoutes(api)
		inventory.RegisterRoutes(api)
		stats.RegisterRoutes(api)
		banners.RegisterRoutes(api, db.DB)
	}

	client.RegisterRoutes(r)
	order.RegisterRoutes(r)
	lombard.RegisterRoutes(r)
	servicerequest.RegisterRoutes(r)
	user.RegisterRoutes(r)
	settings.RegisterRoutes(r)

	return r
}
