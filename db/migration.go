package db

import (
	"fmt"
	"main/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)


func Admin_Migration(db *gorm.DB){


	if db.Migrator().HasTable(&models.Memberships{}) {
		var membership models.Memberships
		
		query := "SELECT memberships.membership_name, memberships.price FROM memberships;"
		db.Raw(query).Scan(&membership)
		if membership.Membership_name=="" {
			Membership_Distribution(db)
		}
	}
	if db.Migrator().HasTable(&models.User{}){

		//create admin
		var admin models.User
		query:="SELECT users.user_id FROM users WHERE role='admin'"
		db.Raw(query).Scan(&admin)
		if admin.User_id==""{
			//create admin 
			Create_Admin(db)
			//give token to this admin
			GetToken(db)
		}

	}
}


func Membership_Distribution(db *gorm.DB){

	//default memberships with their prices
	var membership models.Memberships

	membership.Membership_name="Individual"
	membership.Price=119

	db.Create(&membership)

	membership.Membership_name="Duo"
	membership.Price=149

	db.Create(&membership)

}


func Create_Admin(db *gorm.DB){


	var admin models.User

	admin.Role="admin"
	admin.Name="aman-admin"
	er:=db.Create(&admin).Error
	if er != nil {
		fmt.Println("db error during admin creation")
	}
	
}


func GetToken(db *gorm.DB){

	// jwt authentication token
	
	expirationTime := time.Now().Add(365* 24 * time.Hour)
	fmt.Println("expiration time is: ", expirationTime)

	// check if the user is valid then only create token

	var user models.User
	er:=db.Where("role=?", "admin").First(&user).Error
	if er != nil {
		fmt.Println("db error ",er)
	}
	claims := models.Claims{

		Role:user.Role,
		User_id:user.User_id,
		Active:true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token: ", token)
	tokenString, err := token.SignedString(os.Getenv("Jwt_key"))
	if err != nil {
		fmt.Println("error is :", err)
		
	}
	
	user.Token=tokenString
	user.LoggedIn=true
	Er:=db.Where("role=?", "admin").Updates(&user).Error
	if Er!=nil{

		fmt.Println("db error",Er)
	}

	
}