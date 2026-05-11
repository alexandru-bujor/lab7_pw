package main

import (
	"log"
	"os"

	"MegaMobileBack/internal/category"
	"MegaMobileBack/internal/client"
	"MegaMobileBack/internal/config"
	"MegaMobileBack/internal/inventory"
	"MegaMobileBack/internal/lombard"
	makepkg "MegaMobileBack/internal/make"
	"MegaMobileBack/internal/order"
	"MegaMobileBack/internal/product"
	"MegaMobileBack/internal/service"
	servicecategory "MegaMobileBack/internal/servicecategory"
	"MegaMobileBack/internal/servicerequest"
	"MegaMobileBack/internal/settings"
	"MegaMobileBack/internal/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	config.LoadEnv()

	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = db.AutoMigrate(
		&makepkg.Make{},
		&makepkg.Model{},
		&servicecategory.ServiceCategory{},
		&product.Product{},
		&product.ProductImage{},
		&product.ProductVariant{},
		&category.CategoryTemplate{},
		&category.Category{},
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
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}
