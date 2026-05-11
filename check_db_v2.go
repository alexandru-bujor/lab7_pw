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

	var count int64
	db.Raw("SELECT count(*) FROM information_schema.columns WHERE table_name = 'orders' AND table_schema = ?", name).Scan(&count)
	fmt.Printf("Total columns in orders table: %d\n", count)

	rows, _ := db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'orders' AND table_schema = ?", name).Rows()
	defer rows.Close()
	fmt.Println("All columns:")
	for rows.Next() {
		var col string
		rows.Scan(&col)
		fmt.Printf("\"%s\", ", col)
	}
	fmt.Println()
}
