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

	fmt.Println("Comparing Manager IDs with User IDs...")

	rows, _ := db.Raw("SELECT id, manager_id FROM orders WHERE manager_id IS NOT NULL").Rows()
	defer rows.Close()

	for rows.Next() {
		var oid, mid int
		rows.Scan(&oid, &mid)

		var exists bool
		db.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", mid).Scan(&exists)
		fmt.Printf("Order ID: %d, Manager ID: %d, User exists: %v\n", oid, mid, exists)
	}
}
