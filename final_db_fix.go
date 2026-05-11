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

	fmt.Println("CRITICAL DB FIX START")

	// Check for 'manager_idorage'
	var colExists int
	db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = 'orders' AND table_schema = ? AND column_name = 'manager_idorage'", name).Scan(&colExists)
	if colExists > 0 {
		fmt.Println("Found broken column 'manager_idorage'. Renaming to 'variant_storage'...")
		db.Exec("ALTER TABLE orders CHANGE manager_idorage variant_storage VARCHAR(255)")
	}

	// Check for 'manager_id'
	db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = 'orders' AND table_schema = ? AND column_name = 'manager_id'", name).Scan(&colExists)
	if colExists == 0 {
		fmt.Println("Column 'manager_id' is MISSING. Adding it now...")
		err = db.Exec("ALTER TABLE orders ADD COLUMN manager_id INT").Error
		if err != nil {
			fmt.Println("FAILED to add manager_id:", err)
		} else {
			fmt.Println("Successfully added manager_id.")
		}
	} else {
		fmt.Println("Column 'manager_id' already exists.")
	}

	// Check if 'variant_storage' is missing
	db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = 'orders' AND table_schema = ? AND column_name = 'variant_storage'", name).Scan(&colExists)
	if colExists == 0 {
		fmt.Println("Column 'variant_storage' is MISSING. Adding it now...")
		db.Exec("ALTER TABLE orders ADD COLUMN variant_storage VARCHAR(255)")
	}

	fmt.Println("CRITICAL DB FIX DONE")
}
