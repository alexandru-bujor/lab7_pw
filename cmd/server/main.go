package main

import (
	"log"
	"os"

	"MegaMobileBack/internal/banners"
	"MegaMobileBack/internal/category"
	"MegaMobileBack/internal/client"
	"MegaMobileBack/internal/config"
	"MegaMobileBack/internal/inventory"
	"MegaMobileBack/internal/lombard"
	makepkg "MegaMobileBack/internal/make"
	"MegaMobileBack/internal/order"
	"MegaMobileBack/internal/product"
	"MegaMobileBack/internal/router"
	"MegaMobileBack/internal/service"
	servicecategory "MegaMobileBack/internal/servicecategory"
	"MegaMobileBack/internal/servicerequest"
	"MegaMobileBack/internal/settings"
	"MegaMobileBack/internal/telegram"
	"MegaMobileBack/internal/user"
	"MegaMobileBack/pkg/db"
)

func main() {
	config.LoadEnv()
	db.Connect()

	// Pre-migration step: Fix empty or null slugs to prevent AutoMigrate Unique Constraint errors (Error 1062)
	if err := db.DB.Exec("UPDATE products SET slug = CONCAT('product-', id) WHERE slug = '' OR slug IS NULL;").Error; err != nil {
		log.Printf("Failed to patch empty product slugs before migration: %v", err)
	}

	if err := db.DB.AutoMigrate(
		&makepkg.Make{},
		&makepkg.Model{},
		&servicecategory.ServiceCategory{},
		&product.Product{},
		&product.ProductImage{},
		&category.CategoryTemplate{},
		&service.Service{},
		&client.Client{},
		&client.ClientOrder{},
		&client.ClientOrderItem{},
		&inventory.InventoryItem{},
		&inventory.InventorySection{},
		&order.Order{},
		&lombard.LombardRequest{},
		&servicerequest.ServiceRequest{},
		&user.User{},
		&settings.Setting{},
		&banners.Banner{},
	); err != nil {
		log.Printf("AutoMigrate failed: %v", err)
	}

	var userCount int64
	if err := db.DB.Model(&user.User{}).Count(&userCount).Error; err != nil {
		log.Printf("Failed to check user count: %v", err)
	} else if userCount == 0 {
		userService := user.NewService()
		adminUser, err := userService.Create(user.CreateUserInput{
			FirstName:   "Admin",
			LastName:    "User",
			Email:       "admin@admin.com",
			Password:    "admin123",
			Role:        "admin",
			Permissions: []string{},
		})
		if err != nil {
			log.Printf("Failed to create default admin user: %v", err)
		} else {
			log.Printf("Default admin user created: %s (ID: %d)", adminUser.Email, adminUser.ID)
		}
	}

	if err := telegram.StartBot(); err != nil {
		log.Printf("Failed to start Telegram bot: %v", err)
	}

	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
