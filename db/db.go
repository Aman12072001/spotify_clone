package db

import (
	"fmt"
	"os"

	mod "main/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	db.AutoMigrate(&mod.User{},&mod.AudioFile{},&mod.Playlist{},&mod.Memberships{},mod.Fav_Songs{},mod.Recently_Played{},mod.Album{},mod.Payments{},mod.Blacklisted_tokens{})

	DB = db


	//migrate admin in starting of app (if not already present)
	Admin_Migration(db)

	return nil


}

