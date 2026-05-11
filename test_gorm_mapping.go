package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

type Order struct {
	ID          int
	ManagerName string `gorm:"-"`
}

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

	var o Order
	db.Table("orders").Select("id, 'Test Name' as manager_name").Limit(1).Scan(&o)
	fmt.Printf("Using Scan: ID: %d, ManagerName: '%s'\n", o.ID, o.ManagerName)

	var orders []Order
	db.Table("orders").Select("id, 'Test Name' as manager_name").Limit(1).Find(&orders)
	if len(orders) > 0 {
		fmt.Printf("Using Find: ID: %d, ManagerName: '%s'\n", orders[0].ID, orders[0].ManagerName)
	}
}
