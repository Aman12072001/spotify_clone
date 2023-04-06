package cont

import (
	"encoding/json"
	"fmt"
	Res "main/Response"
	"main/db"
	"main/models"
	con "main/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)





func User_login_with_contact_no(w http.ResponseWriter,r *http.Request){

	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)

		return
	}  
	
	w.Header().Set("Content-Type", "application/json")
	
	//take the user name ,email and contact number

	var user models.User
	input :=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	contact:=input["contact"]
	if len([]rune(contact))>10 || len([]rune(contact))==0   {
		Res.Response("Bad request ",400,"enter valid contact number of 10 digits","",w)

	}
	if _, err := strconv.ParseInt(contact,10,64); err == nil {
		
	}

	json.NewDecoder(r.Body).Decode(&user)

	er:=db.DB.Create(&user).Error
	if er!=nil{
		Res.Response("server error",500,er.Error(),"",w)

	}

	//send otp according to the contact number entered
	sendOtp("+91" + user.Contact_no)
	//generate an Otp
	Res.Response("OK",200,"OTP SENT","",w)


}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}



func UpdateProfile(w http.ResponseWriter,r *http.Request){

//update user information facilities
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)

	}

	w.Header().Set("content-type", "application/json")
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	//convert the user password to hash
	
	if user.Password!=""{

		pwd,_:=HashPassword(user.Password)

		user.Password=pwd

	}
	

	
	// func CheckPasswordHash(password, hash string) bool {
	// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// 	return err == nil
	// }

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
		if claims.Role!="admin" {
			claims.Role=""//do not allow user to change itsa role to admin
		}
		if claims.User_id==user.User_id{

			//update information now
			err:=db.DB.Where("user_id=?",user.User_id).Updates(&user).Error
			if err != nil {
				fmt.Println("err",err.Error())
				Res.Response("server error",500,err.Error(),"",w)


			}
			//fmt.Fprint(w,"Profile updated successfully")
			Res.Response("OK",200,"Profile updated successfully","",w)


		}
		
		
	} else {
		fmt.Println(err)
		Res.Response("Unauthorized",401,err.Error(),"",w)

	}




}




func GetSong(w http.ResponseWriter,r * http.Request){


	// //get the song from db based on the id of song
	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")

	var song models.AudioFile

	json.NewDecoder(r.Body).Decode(&song)

	err:=db.DB.Where("id=?",song.ID).First(&song).Error
	if err!=nil{

		Res.Response("server error",500,err.Error(),"",w)

	}



		//return the path for the frontend dev

		// response.Code=200
		// response.Status="OK"
		// response.Message="token provided successfully"
	// json.NewEncoder(w).Encode(&song)
	Res.Response("OK",200,"Success",song,w)


}

func Get_All_Songs(w http.ResponseWriter,r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	var songs []models.AudioFile

	query:="SELECT * FROM audio_files;"
	err:=db.DB.Raw(query).Scan(&songs).Error
	if err!=nil{
		Res.Response("server error",500,err.Error(),"",w)

	}
	Res.Response("OK",200,"Success",songs,w)

	// json.NewEncoder(w).Encode(&songs)



}

func CreatePlaylist(w http.ResponseWriter,r * http.Request){


	//custom playlist
	//user want to add songs to his/her playlist

	//it will take  playlist_name and path 
	//user_id will be set from the token
	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}
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
		Res.Response("Unauthorized",401,"token not valid","",w)
	}

	db.DB.Create(&playlist)

	// fmt.Fprint(w,"added to playlist")
	Res.Response("OK",200,"added to playlist","",w)


}

func Show_playlist(w http.ResponseWriter, r *http.Request){


	//take the name of the playlist
	//userid will be automatically fetch from token
	w.Header().Set("Content-Type", "application/json")

	var playlist models.Playlist //for decoding body
	var playlists_song []models.Playlist //for display all songs of playlist

	json.NewDecoder(r.Body).Decode(&playlist)
	// fmt.Println("playlist name",playlist.Playlist_name)

	//token parsing for user credentials
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
		Res.Response("Unauthorized",401,"token not valid","",w)

	}

	
	
	
	query:="SELECT * FROM playlists WHERE playlist_name='"+ playlist.Playlist_name +"'AND user_id='"+ playlist.User_id+"';"

	//query_with_pagination:="SELECT *FROM playlists WHERE playlist_name='"+ playlist.Playlist_name +"'LIMIT ;"
	// fmt.Println("query : ",query)

	er:=db.DB.Raw(query).Scan(&playlists_song).Error
	if er!=nil{
		Res.Response("server error",500,er.Error(),"",w)
	}
	// json.NewEncoder(w).Encode(&playlists_song)
	Res.Response("OK",200,"Success",playlists_song,w)
	


}

func Add_song_toFav(w http.ResponseWriter,r *http.Request){


	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}

	w.Header().Set("Content-Type", "application/json")
	//take the song path as input from r.body
	var song models.AudioFile

	json.NewDecoder(r.Body).Decode(&song)

	var fav_song models.Fav_Songs

	//parse the token to get the user_id
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
		fav_song.User_id=claims.User_id //user id milgyi
	} else {
		fmt.Println(err)
		Res.Response("Unauthorized",401,"token not valid","",w)
	}

	
	fav_song.Song_id=song.ID//song ka path set hogya


	er:=db.DB.Create(&fav_song).Error // table me create hogya
	if er != nil {

		Res.Response("server error",500,er.Error(),"",w)
	}



}


func Add_to_RecentlyPlayed(w http.ResponseWriter,r *http.Request){


	//userid will be fetch from token
	//playedAt time will be set manually
	//song id will be decoded from r.body
	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}

	var recent_add_song models.Recently_Played

	json.NewDecoder(r.Body).Decode(&recent_add_song)


		recent_add_song.PlayedAt=time.Now() //played at this time 
		

		//parse the token to get userid
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
			recent_add_song.User_id=claims.User_id //user id milgyi
		} else {

			fmt.Println(err)
			Res.Response("Unauthorized",401,"token not valid","",w)
		}


		er:=db.DB.Create(&recent_add_song).Error
		if er!=nil{
			Res.Response("server error",500,er.Error(),"",w)
		}
		// fmt.Fprint(w,"added to recently played")
		Res.Response("OK",200,"token provided successfully","",w)
		

	
}


func Get_Recently_Played_Songs(w http.ResponseWriter,r *http.Request){

	//get the list of 20 recently played songs
	w.Header().Set("Content-Type", "application/json")

	//get the userid from token parsing-->(to know who is asking for this request)
	//filter the entries of recently played based on userid

	var user models.User

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
		//user id milgyi
		user.User_id=claims.User_id
	} else {
		fmt.Println(err)
		Res.Response("Unauthorized",401,"token not valid","",w)
	}


	query:="SELECT * FROM recently_playeds WHERE user_id='"+user.User_id+"' ORDER BY  played_at DESC LIMIT 20"
	var list_of_recently_played []models.Recently_Played
	er:=db.DB.Raw(query).Scan(&list_of_recently_played).Error
	if er!=nil{

		Res.Response("server error",500,er.Error(),"",w)
	}

	Res.Response("OK",200,"success",list_of_recently_played,w)
	// json.NewEncoder(w).Encode(&list_of_recently_played)



}

func Get_Artist(w http.ResponseWriter,r *http.Request){


	//get the name of the artist from r.body
	w.Header().Set("Content-Type", "application/json")


	var artist models.Artist

	json.NewDecoder(r.Body).Decode(&artist)

	//based on the artist name get the songs from audiofiles table

	query:="SELECT * FROM audio_files WHERE artist='"+artist.Artist_name+"';"

	var list_of_artist_songs []models.AudioFile

	er:=db.DB.Raw(query).Scan(&list_of_artist_songs).Error
	if er!=nil{

		Res.Response("server error",500,er.Error(),"",w)
	}

	// json.NewEncoder(w).Encode(&list_of_artist_songs)
	Res.Response("OK",200,"Success",list_of_artist_songs,w)
}

func Get_Album(w http.ResponseWriter,r *http.Request){

		//get album based on the name of album
		if r.Method != http.MethodPost {
			// w.WriteHeader(http.StatusMethodNotAllowed)
			Res.Response("Method Not Allowed ",405,"use correct http method","",w)
			
		}
		w.Header().Set("Content-Type", "application/json")

		var album models.Album
		var album_song_list []models.Album
		json.NewDecoder(r.Body).Decode(&album)

		query:="SELECT * FROM albums WHERE album_name='"+album.Album_name+"';"
		er:=db.DB.Raw(query).Scan(&album_song_list).Error
		if er!=nil{

			Res.Response("server error",500,er.Error(),"",w)
		}
		Res.Response("OK",200,"Success",album_song_list,w)
		// json.NewEncoder(w).Encode(&album_song_list)
	
}