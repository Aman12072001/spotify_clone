package cont

import (
	"encoding/json"
	"fmt"
	"main/db"
	"main/models"
	con "main/utils"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v4"
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

func UpdateProfile(w http.ResponseWriter,r *http.Request){

//update user information facilities

	w.Header().Set("content-type", "application/json")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	//check whether user has the correct token to change user information
	parsedToken ,err := jwt.ParseWithClaims(r.Header["Token"][0] ,&models.Claims{}, func(token *jwt.Token) (interface{}, error) {
						
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("error")
		}
		return con.Jwt_key , nil
	})

	fmt.Println("token parsing hogyi")

	if claims, ok := parsedToken.Claims.(*models.Claims); ok && parsedToken.Valid {
		// fmt.Printf("token will expire at :%v",  claims.ExpiresAt)
		// fmt.Println("claims ki userid",claims)
		if claims.User_id==user.User_id{

			//update information now
			err:=db.DB.Where("user_id=?",user.User_id).Updates(&user).Error
			if err != nil {
				fmt.Println("err",err.Error())

			}
			fmt.Fprint(w,"Profile updated successfully")

		}
		
		
	} else {
		fmt.Println(err)
	}




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

func Get_All_Songs(w http.ResponseWriter,r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	var songs []models.AudioFile

	query:="SELECT * FROM audio_files;"
	db.DB.Raw(query).Scan(&songs)

	json.NewEncoder(w).Encode(&songs)



}

func CreatePlaylist(w http.ResponseWriter,r * http.Request){


	//custom playlist
	//user want to add songs to his/her playlist

	//it will take  playlist_name and path 
	//user_id will be set from the token
	w.Header().Set("Content-Type", "application/json")
	var playlist models.Playlist


	json.NewDecoder(r.Body).Decode(&playlist)

	//extract the user_id from the token
	fmt.Println("playlist var me value encode ho gyi")
	fmt.Println("header token vlaue",r.Header["Token"][0])

	parsedToken ,err := jwt.ParseWithClaims(r.Header["Token"][0] ,&models.Claims{}, func(token *jwt.Token) (interface{}, error) {
						
		if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("error")
		}
		return con.Jwt_key , nil
	})

	fmt.Println("token parsing hogyi")

	if claims, ok := parsedToken.Claims.(*models.Claims); ok && parsedToken.Valid {
		// fmt.Printf("token will expire at :%v",  claims.ExpiresAt)
		fmt.Println("claims ki userid",claims)
		playlist.User_id=claims.User_id
	} else {
		fmt.Println(err)
	}

	db.DB.Create(&playlist)

	fmt.Fprint(w,"added to playlist")


}

func Show_playlist(w http.ResponseWriter, r *http.Request){


	//take the name of the playlist
	//userid will be automatically fetch from token
	w.Header().Set("Content-Type", "application/json")

	var playlist models.Playlist //for decoding body
	var playlists []models.Playlist //for display all songs of playlist

	json.NewDecoder(r.Body).Decode(&playlist)
	// fmt.Println("playlist name",playlist.Playlist_name)

	
	
	
	query:="SELECT * FROM playlists WHERE playlist_name='"+ playlist.Playlist_name +"';"

	//query_with_pagination:="SELECT *FROM playlists WHERE playlist_name='"+ playlist.Playlist_name +"'LIMIT ;"
	// fmt.Println("query : ",query)

	db.DB.Raw(query).Scan(&playlists)

	json.NewEncoder(w).Encode(&playlists)
	


}