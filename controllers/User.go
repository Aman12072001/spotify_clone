package cont

import (
	"encoding/json"
	"main/db"
	"main/models"
	"net/http"
	"os"
)


func User_login(w http.ResponseWriter,r *http.Request){

	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	
	//take the user name ,email and contact number

	var user models.User

	json.NewDecoder(r.Body).Decode(&user)

	db.DB.Create(&user)

	//send otp according to the contact number entered
	sendOtp("+91" + user.Contact_no)
	//generate an Otp


}





func GetSong(w http.ResponseWriter,r * http.Request){


	// //get the song from db based on the name of song
	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")

	var song models.AudioFile

	json.NewDecoder(r.Body).Decode(&song)

	db.DB.Where("name=?",song.Name).First(&song)

	//set the priviledge for this file ,so that frontend can access it
	// set file permissions to read and write for owner, read only for group and others
		err := os.Chmod(song.Path, 0750)
		if err != nil {
			panic(err)
		}

		//return the path for the frontend dev

	

	

	json.NewEncoder(w).Encode(&song)




}

func CreatePlaylist()