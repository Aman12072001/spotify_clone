package models

import (
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
	
}

type Playlist struct {

	User_id   string `json:"user_id"`
	Playlist_name string `json:"playlist_name"`
	Song_Path  string `json:"song_path"`
	

}

type AudioFile struct {
    ID        string   `json:"id" gorm:"default:uuid_generate_v4();"`
	Path      string `json:"path" gorm:"unique"`
    Name      string `json:"name"`
    Size      float64`json:"size"`
	Artist    string `json:"artist"`
    
}

type Paths struct{

	Path  string `json:"path"`
}

type Claims struct{

	Phone string `json:"phone"`
	User_id string `json:"user_id"`
	Role string `json:"role"`
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

