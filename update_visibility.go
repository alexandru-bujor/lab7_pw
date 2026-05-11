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

	fmt.Println("Setting is_visible_in_secondary = true for manual orders...")
	res := db.Exec("UPDATE orders SET is_visible_in_secondary = true WHERE order_source = 'manual'")
	if res.Error != nil {
		fmt.Printf("ERROR: %v\n", res.Error)
	} else {
		fmt.Printf("Success! Rows affected: %d\n", res.RowsAffected)
	}
}
