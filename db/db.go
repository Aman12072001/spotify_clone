package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	mod "main/models"
)
var DB *gorm.DB
func Connect()error{


	fmt.Println("Connecting to database...")

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", os.Getenv("DB_Host"),os.Getenv("DB_Port"), os.Getenv("DB_User"), os.Getenv("DB_Password"), os.Getenv("Dbname"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Error in connecting to database:", err)
		return err
	}

	db.AutoMigrate(&mod.User{},&mod.AudioFile{})

	DB = db

	return nil


}