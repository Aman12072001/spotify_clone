package models

import(
	jwt "github.com/golang-jwt/jwt/v4"
)

type User struct{

	User_id        string   `json:"user_id" gorm:"default:uuid_generate_v4();"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	// Playlist []string `json:"playlist"`
	Membership string `json:"membership"`
	Contact_no string `json:"contact gorm:default:unique" `
	Token  string  `json:"token"`
	
}

type Playlist struct {

	User_id   string `json:"user_id"`
	Song_name string  `json:"song_name"`
	

}

type AudioFile struct {
    ID        string   `json:"id" gorm:"default:uuid_generate_v4();"`
	Path      string `json:"path"`
    Name      string `json:"name"`
    Size      float64`json:"size"`
    
}

type Paths struct{

	Path  string `json:"path"`
}

type Claims struct{

	Phone string `json:"phone"`
	jwt.RegisteredClaims
		
	
}