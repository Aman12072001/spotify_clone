package models

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type User struct{

	User_id        string   `json:"user_id" gorm:"default:uuid_generate_v4();"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	// Playlist []string `json:"playlist"`
	Membership string `json:"membership"`
	Contact_no string `json:"contact" `
	Token  string  `json:"token"`
	Role string `json:"role"`
	LoggedIn bool `json:"logged_in"`
	
}

type Playlist struct {

	User_id   string `json:"user_id"`
	Playlist_name string `json:"playlist_name"`
	Song_id  string `json:"song_id"`
	

}

type AudioFile struct {
    ID        string   `json:"id" gorm:"default:uuid_generate_v4();"`
	Path      string `json:"path" gorm:"unique"`
    Name      string `json:"name"`
    Size      float64`json:"size"`
	Artist    string `json:"artist"`
	Img_Path string `json:"img_path"`
    
}

type Path struct {

	Path   string `json:"path"`
}

type Claims struct{

	Phone string `json:"phone"`
	User_id string `json:"user_id"`
	Role string `json:"role"`
	Active bool `json:"active"`
	jwt.RegisteredClaims
		
}

type Memberships struct{

	Membership_name string `json:"membership_name"`
	Price float64 `json:"price"`

}

type Artist struct{


	Artist_name string `json:"artist_name"`
	Artist_song string `json:"artist_song"`
}

type Fav_Songs struct{

	User_id string 
	Song_id string 


}

type Recently_Played struct{


	User_id string 
	Song_id string  `json:"song_id"`
	PlayedAt time.Time
}

type Album struct{

	Album_name string `json:"album_name"`
	// Artist_name string `json:"artist_name"`
	Song_id string `json:"song_id"`
}


type Payments struct{

	User_id string `json:"user_id"`
	Order_id string `json:"order_id"`
	Payment_id string `json:"payment_id"`
	Membership_name string `json:"membershipname"`
	Time time.Time `json:"time"`
	Status string `json:"status"`
}


type Response struct {
	Status  string      `json:"status"`
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
