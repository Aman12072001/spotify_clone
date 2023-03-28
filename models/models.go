package models


type User struct{

	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	Playlist []byte `json:"playlist"`
	Membership string `json:"membership"`
	Contact_no string `json:"contact"`
	
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