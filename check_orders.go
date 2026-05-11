package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Order struct {
	ID        int
	ManagerID *int
	Status    string
	CreatedAt time.Time
}

func main() {
	godotenv.Load(".env")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var orders []Order
	db.Table("orders").Order("id desc").Limit(5).Find(&orders)

	fmt.Println("Recent Orders:")
	for _, o := range orders {
		mid := "nil"
		if o.ManagerID != nil {
			mid = fmt.Sprintf("%d", *o.ManagerID)
		}
		fmt.Printf("ID: %d, Status: %s, ManagerID: %s, CreatedAt: %v\n", o.ID, o.Status, mid, o.CreatedAt)
	}
}
