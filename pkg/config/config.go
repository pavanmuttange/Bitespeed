package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", "dpg-cpohueij1k6c73a7ku50-a", "pavandb_user", "VpUJWoN5TYNlPRBoMTrWij0jWHYqoivS", "pavandb", "5432")
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", "localhost", "postgres", "postgres", "postgres", "5432")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("error in connecting database")
		return
	}

	fmt.Println("Database connected successfully")

}

func Init() {
	ConnectDB()
}
