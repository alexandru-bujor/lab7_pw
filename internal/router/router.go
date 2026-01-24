package router

import (
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

	"log"
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

	// Log incoming Origin header for debugging CORS issues
	r.Use(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Writer.Header().Add("X-Debug-Origin-Received", origin)
			log.Printf("Incoming request %s %s Origin=%s", c.Request.Method, c.Request.URL.Path, origin)
		}
		c.Next()
	})

	// Disable automatic trailing slash redirect (causes CORS issues)
	r.RedirectTrailingSlash = false

	// CORS: Support both development and production origins
	// MUST be before other middleware to handle preflight requests
	allowedOrigins := []string{
		"http://localhost:8081",
		"http://localhost:5173",
		"http://localhost:8080",
		"http://localhost:3000",
		"https://megamobile.md",
		"https://www.megamobile.md",
	}

	// Allow additional origins from environment variable (comma-separated)
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

	// Health endpoint
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

	// Serve static files from uploads directory
	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		// auth
		api.POST("/login", user.LoginHandler)

		// products
		productsGroup := api.Group("/products")
		product.Routes(productsGroup)

		// categories
		categoriesGroup := api.Group("/categories")
		category.Routes(categoriesGroup)

		// services
		service.RegisterRoutes(r)

		// service categories
		servicecategory.RegisterRoutes(r)

		// makes and models
		make.RegisterRoutes(r)

		// upload
		upload.RegisterRoutes(api)

		// inventory
		inventory.RegisterRoutes(api)

		// stats
		stats.RegisterRoutes(api)
	}

	// clients
	client.RegisterRoutes(r)

	// orders
	order.RegisterRoutes(r)

	// lombard
	lombard.RegisterRoutes(r)

	// service requests
	servicerequest.RegisterRoutes(r)

	// users
	user.RegisterRoutes(r)

	// settings
	settings.RegisterRoutes(r)

	return r
}
