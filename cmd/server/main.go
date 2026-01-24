package main

import (
	"log"
	"os"

	"MegaMobileBack/internal/config"
	"MegaMobileBack/internal/product"
	"MegaMobileBack/internal/router"
	"MegaMobileBack/internal/category"
	"MegaMobileBack/internal/service"
	"MegaMobileBack/internal/client"
	"MegaMobileBack/internal/inventory"
	"MegaMobileBack/internal/order"
	"MegaMobileBack/pkg/db"
	"MegaMobileBack/internal/lombard"
	"MegaMobileBack/internal/servicerequest"
	"MegaMobileBack/internal/user"
	"MegaMobileBack/internal/telegram"
	servicecategory "MegaMobileBack/internal/servicecategory"
	makepkg "MegaMobileBack/internal/make"
	"MegaMobileBack/internal/settings"
)

func main() {
	       config.LoadEnv()
	       db.Connect()

	       migrator := db.DB.Migrator()

	       // Create makes and models tables if missing
	       if err := migrator.AutoMigrate(&makepkg.Make{}, &makepkg.Model{}); err != nil {
	       log.Printf("⚠️ AutoMigrate makes/models failed: %v", err)
	       } else {
	       log.Println("✅ makes/models tables ready")
	       }

	       // Create service_categories table if missing (after migrator is defined)
	       if err := migrator.AutoMigrate(&servicecategory.ServiceCategory{}); err != nil {
		       log.Printf("⚠️ AutoMigrate service_categories failed: %v", err)
	       } else {
		       log.Println("✅ service_categories table ready")
	       }

	// Create products table and related tables if missing (must be done before adding columns)
	log.Println("🔄 Attempting to create products table...")
	if err := migrator.AutoMigrate(&product.Product{}, &product.ProductImage{}, &product.ProductVariant{}); err != nil {
		log.Printf("❌ AutoMigrate products tables failed: %v", err)
		log.Printf("⚠️ Skipping column additions for products table")
	} else {
		// Verify table was actually created
		if migrator.HasTable(&product.Product{}) {
			log.Println("✅ products, product_images, product_variants tables ready")
		} else {
			log.Printf("❌ AutoMigrate succeeded but products table does not exist!")
			log.Printf("⚠️ Skipping column additions for products table")
		}
		
		// Only add columns if table exists
		if migrator.HasTable(&product.Product{}) {
			// Add 'specs' JSON column if missing
			if !migrator.HasColumn(&product.Product{}, "Specs") {
				if err := migrator.AddColumn(&product.Product{}, "Specs"); err != nil {
					log.Printf("⚠️ AddColumn specs failed: %v", err)
				} else {
					log.Println("✅ Added 'specs' column to products")
				}
			}
			// Add 'condition' column if missing
			if !migrator.HasColumn(&product.Product{}, "Condition") {
				if err := migrator.AddColumn(&product.Product{}, "Condition"); err != nil {
					log.Printf("⚠️ AddColumn condition failed: %v", err)
				} else {
					log.Println("✅ Added 'condition' column to products")
				}
			}
			// Add 'slug' column if missing
			if !migrator.HasColumn(&product.Product{}, "Slug") {
				if err := migrator.AddColumn(&product.Product{}, "Slug"); err != nil {
					log.Printf("⚠️ AddColumn slug failed: %v", err)
				} else {
					log.Println("✅ Added 'slug' column to products")
				}
			}
			// Add brand_name column if missing
			if !migrator.HasColumn(&product.Product{}, "BrandName") {
				if err := migrator.AddColumn(&product.Product{}, "BrandName"); err != nil {
					log.Printf("⚠️ AddColumn brand_name failed: %v", err)
				} else {
					log.Println("✅ Added 'brand_name' column to products")
				}
			}
			// Add variant columns if missing
			if !migrator.HasColumn(&product.Product{}, "VariantGroupID") {
				if err := migrator.AddColumn(&product.Product{}, "VariantGroupID"); err != nil {
					log.Printf("⚠️ AddColumn variant_group_id failed: %v", err)
				} else {
					log.Println("✅ Added 'variant_group_id' column to products")
				}
			}
			if !migrator.HasColumn(&product.Product{}, "VariantColor") {
				if err := migrator.AddColumn(&product.Product{}, "VariantColor"); err != nil {
					log.Printf("⚠️ AddColumn variant_color failed: %v", err)
				} else {
					log.Println("✅ Added 'variant_color' column to products")
				}
			}
			if !migrator.HasColumn(&product.Product{}, "VariantStorage") {
				if err := migrator.AddColumn(&product.Product{}, "VariantStorage"); err != nil {
					log.Printf("⚠️ AddColumn variant_storage failed: %v", err)
				} else {
					log.Println("✅ Added 'variant_storage' column to products")
				}
			}
			if !migrator.HasColumn(&product.Product{}, "VariantRAM") {
				if err := migrator.AddColumn(&product.Product{}, "VariantRAM"); err != nil {
					log.Printf("⚠️ AddColumn variant_ram failed: %v", err)
				} else {
					log.Println("✅ Added 'variant_ram' column to products")
				}
			}
			// Add accounting_type column to products if missing
			if !migrator.HasColumn(&product.Product{}, "AccountingType") {
				if err := migrator.AddColumn(&product.Product{}, "AccountingType"); err != nil {
					log.Printf("⚠️ AddColumn accounting_type failed: %v", err)
				} else {
					log.Println("✅ Added 'accounting_type' column to products")
				}
			}
		} else {
			log.Printf("⚠️ products table does not exist, skipping column additions")
		}
	}
	// Create category_templates table if missing
	if err := migrator.AutoMigrate(&category.CategoryTemplate{}); err != nil {
		log.Printf("⚠️ AutoMigrate category_templates failed: %v", err)
	} else {
		log.Println("✅ category_templates ready")
	}

	// Create services table if missing
	if err := migrator.AutoMigrate(&service.Service{}); err != nil {
		log.Printf("⚠️ AutoMigrate services failed: %v", err)
	} else {
		log.Println("✅ services table ready")
	}
	
	// Add 'prices' JSON column if missing
	if !migrator.HasColumn(&service.Service{}, "Prices") {
		if err := migrator.AddColumn(&service.Service{}, "Prices"); err != nil {
			log.Printf("⚠️ AddColumn prices failed: %v", err)
		} else {
			log.Println("✅ Added 'prices' column to services")
		}
	}
	
	// Add 'photo_url' column if missing
	if !migrator.HasColumn(&service.Service{}, "PhotoURL") {
		if err := migrator.AddColumn(&service.Service{}, "PhotoURL"); err != nil {
			log.Printf("⚠️ AddColumn photo_url failed: %v", err)
		} else {
			log.Println("✅ Added 'photo_url' column to services")
		}
	}

	// Create clients tables
	if err := migrator.AutoMigrate(&client.Client{}, &client.ClientOrder{}, &client.ClientOrderItem{}); err != nil {
		log.Printf("⚠️ AutoMigrate clients tables failed: %v", err)
	} else {
		log.Println("✅ clients, client_orders, client_order_items tables ready")
	}

	// Create inventory table
	if err := migrator.AutoMigrate(&inventory.InventoryItem{}); err != nil {
		log.Printf("⚠️ AutoMigrate inventory_items failed: %v", err)
	} else {
		log.Println("✅ inventory_items table ready")
	}

	// Create inventory_sections table
	if err := migrator.AutoMigrate(&inventory.InventorySection{}); err != nil {
		log.Printf("⚠️ AutoMigrate inventory_sections failed: %v", err)
	} else {
		log.Println("✅ inventory_sections table ready")
	}

	// Create orders table
	if err := migrator.AutoMigrate(&order.Order{}); err != nil {
		log.Printf("⚠️ AutoMigrate orders failed: %v", err)
	} else {
		log.Println("✅ orders table ready")
	}
	// Add accounting_type column to orders if missing
	if !migrator.HasColumn(&order.Order{}, "AccountingType") {
		if err := migrator.AddColumn(&order.Order{}, "AccountingType"); err != nil {
			log.Printf("⚠️ AddColumn accounting_type to orders failed: %v", err)
		} else {
			log.Println("✅ Added 'accounting_type' column to orders")
		}
	}

	// Create lombard_requests table
	if err := migrator.AutoMigrate(&lombard.LombardRequest{}); err != nil {
		log.Printf("⚠️ AutoMigrate lombard_requests failed: %v", err)
	} else {
		log.Println("✅ lombard_requests table ready")
	}

	// Create service_requests table
	if err := migrator.AutoMigrate(&servicerequest.ServiceRequest{}); err != nil {
		log.Printf("⚠️ AutoMigrate service_requests failed: %v", err)
	} else {
		log.Println("✅ service_requests table ready")
	}

	// Create users table and add permissions column if missing
	if err := migrator.AutoMigrate(&user.User{}); err != nil {
		log.Printf("⚠️ AutoMigrate users failed: %v", err)
	} else {
		log.Println("✅ users table ready")
	}
	// Add permissions column if missing
	if !migrator.HasColumn(&user.User{}, "Permissions") {
		if err := migrator.AddColumn(&user.User{}, "Permissions"); err != nil {
			log.Printf("⚠️ AddColumn permissions failed: %v", err)
		} else {
			log.Println("✅ Added 'permissions' column to users")
		}
	}

	// Create default admin user if no users exist
	var userCount int64
	if err := db.DB.Model(&user.User{}).Count(&userCount).Error; err != nil {
		log.Printf("⚠️ Failed to check user count: %v", err)
	} else if userCount == 0 {
		log.Println("🔄 No users found, creating default admin user...")
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
			log.Printf("❌ Failed to create default admin user: %v", err)
		} else {
			log.Printf("✅ Default admin user created: %s (ID: %d)", adminUser.Email, adminUser.ID)
		}
	} else {
		log.Printf("ℹ️ Found %d existing user(s), skipping default admin creation", userCount)
	}

	// Create settings table if missing
	if err := migrator.AutoMigrate(&settings.Setting{}); err != nil {
		log.Printf("⚠️ AutoMigrate settings failed: %v", err)
	} else {
		log.Println("✅ settings table ready")
	}

	// Add is_visible_in_secondary columns to all tables
	log.Println("🔄 Adding is_visible_in_secondary columns...")
	
	// Products
	if !migrator.HasColumn(&product.Product{}, "IsVisibleInSecondary") {
		if err := migrator.AddColumn(&product.Product{}, "IsVisibleInSecondary"); err != nil {
			log.Printf("⚠️ AddColumn is_visible_in_secondary to products failed: %v", err)
		} else {
			log.Println("✅ Added 'is_visible_in_secondary' column to products")
		}
	}
	
	// Orders
	if !migrator.HasColumn(&order.Order{}, "IsVisibleInSecondary") {
		if err := migrator.AddColumn(&order.Order{}, "IsVisibleInSecondary"); err != nil {
			log.Printf("⚠️ AddColumn is_visible_in_secondary to orders failed: %v", err)
		} else {
			log.Println("✅ Added 'is_visible_in_secondary' column to orders")
		}
	}
	
	// Clients
	if !migrator.HasColumn(&client.Client{}, "IsVisibleInSecondary") {
		if err := migrator.AddColumn(&client.Client{}, "IsVisibleInSecondary"); err != nil {
			log.Printf("⚠️ AddColumn is_visible_in_secondary to clients failed: %v", err)
		} else {
			log.Println("✅ Added 'is_visible_in_secondary' column to clients")
		}
	}
	
	// Inventory
	if !migrator.HasColumn(&inventory.InventoryItem{}, "IsVisibleInSecondary") {
		if err := migrator.AddColumn(&inventory.InventoryItem{}, "IsVisibleInSecondary"); err != nil {
			log.Printf("⚠️ AddColumn is_visible_in_secondary to inventory_items failed: %v", err)
		} else {
			log.Println("✅ Added 'is_visible_in_secondary' column to inventory_items")
		}
	}
	
	// Services
	if !migrator.HasColumn(&service.Service{}, "IsVisibleInSecondary") {
		if err := migrator.AddColumn(&service.Service{}, "IsVisibleInSecondary"); err != nil {
			log.Printf("⚠️ AddColumn is_visible_in_secondary to services failed: %v", err)
		} else {
			log.Println("✅ Added 'is_visible_in_secondary' column to services")
		}
	}
	
	// Lombard requests
	if !migrator.HasColumn(&lombard.LombardRequest{}, "IsVisibleInSecondary") {
		if err := migrator.AddColumn(&lombard.LombardRequest{}, "IsVisibleInSecondary"); err != nil {
			log.Printf("⚠️ AddColumn is_visible_in_secondary to lombard_requests failed: %v", err)
		} else {
			log.Println("✅ Added 'is_visible_in_secondary' column to lombard_requests")
		}
	}

	// Start Telegram bot
	if err := telegram.StartBot(); err != nil {
		log.Printf("⚠️ Failed to start Telegram bot: %v", err)
		log.Println("⚠️ Continuing without Telegram bot...")
	}

	r := router.Setup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
