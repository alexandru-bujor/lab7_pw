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

	fmt.Println("Attempting to update an order with manager_id...")
	res := db.Exec("UPDATE orders SET manager_id = 3 WHERE id = (SELECT id FROM (SELECT id FROM orders LIMIT 1) as t)")
	if res.Error != nil {
		fmt.Printf("ERROR: %v\n", res.Error)
	} else {
		fmt.Printf("Success! Rows affected: %d\n", res.RowsAffected)
	}
}
