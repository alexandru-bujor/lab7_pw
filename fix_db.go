package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	godotenv.Load(".env")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// 1. Check if manager_idorage exists and rename it back to variant_storage
	fmt.Println("Fixing columns...")

	// Check if manager_idorage exists
	var count int64
	db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = 'orders' AND table_schema = ? AND column_name = 'manager_idorage'", name).Scan(&count)
	if count > 0 {
		fmt.Println("Renaming manager_idorage to variant_storage...")
		err = db.Exec("ALTER TABLE orders CHANGE manager_idorage variant_storage VARCHAR(255)").Error
		if err != nil {
			fmt.Println("Error renaming column:", err)
		}
	} else {
		fmt.Println("manager_idorage not found.")
	}

	// 2. Add manager_id column if it doesn't exist
	db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = 'orders' AND table_schema = ? AND column_name = 'manager_id'", name).Scan(&count)
	if count == 0 {
		fmt.Println("Adding manager_id column...")
		err = db.Exec("ALTER TABLE orders ADD COLUMN manager_id INT").Error
		if err != nil {
			fmt.Println("Error adding manager_id:", err)
		} else {
			// Add foreign key
			fmt.Println("Adding foreign key...")
			err = db.Exec("ALTER TABLE orders ADD CONSTRAINT fk_orders_manager FOREIGN KEY (manager_id) REFERENCES users(id)").Error
			if err != nil {
				fmt.Println("Error adding FK:", err)
			}
		}
	} else {
		fmt.Println("manager_id column already exists.")
	}

	fmt.Println("Done.")
}
