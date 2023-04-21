package cont

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	Res "main/Response"
	"main/db"
	"main/models"
	"main/utils"
	con "main/utils"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
)



func User_SignUp(w http.ResponseWriter,r *http.Request){

	utils.SetHeader(w)
	
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)

		return
	}  
	

	
	input :=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	
	err := validation.Validate(input,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("email", validation.Required, is.Email,),
			
			validation.Key("create_password",validation.Required,validation.Length(8, 16),),

			validation.Key("contact",validation.Required),
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}

	//if email already exists 
	query:="select exists(select * from users where email='"+input["email"]+"');"
	var user_exists bool
	db.DB.Raw(query).Scan(&user_exists)
	if user_exists{

		Res.Response("Bad Request",400,"email already exists","",w)

	}else{
		//create db entry for signup input
		var user models.User
		user.Email=input["email"]
		user.Password,_=HashPassword(input["create_password"])
		user.Contact_no=input["contact"]
		er:=db.DB.Create(&user).Error
		if er!=nil{
			Res.Response("Server error",500,er.Error(),"",w)
			return
		}
		Res.Response("Success",200,"user sign up complete","",w)


	}

	


	//redirect to login page 


}


func Forgot_Password(w http.ResponseWriter, r *http.Request){


	utils.SetHeader(w)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)

		return
	} 


	input :=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)

	err := validation.Validate(input,
		validation.Map(
			//contact number must be between 8-10 digits.
			validation.Key("email",validation.Required,is.Email),
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}

	//check whether the email exists in the database
	query:="select exists(select * from users where email='"+input["email"]+"');"

	var emailExists bool

	db.DB.Raw(query).Scan(&emailExists)

	if emailExists{

			//give cookie with email in claims
			claims :=models.Claims{ Role:"forgotuser",Email: input["email"]}
			forgotuserToken:=GenerateNewToken(&claims)
			cookie:=&http.Cookie{Name:"forgotuser_cookie",Value:forgotuserToken}
			http.SetCookie(w,cookie)
			fmt.Println("claims email",claims.Email)

		//send the otp or link on email
		SendEmail(input["email"])
		Res.Response("Sucess",200,"Email Sent","",w)
		

		fmt.Println("cookie forgotuser is set")
		//----------------------------------//

	}else{

		Res.Response("Bad Request",400,"Unknown email address","",w)
		return
	}





}

func SendEmail(toEmail string){


	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "amantarar01@gmail.com")
  
	// Set E-Mail receivers
	m.SetHeader("To", toEmail)
  
	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")
  
	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", ` <html>
  
	<body>
	
	<button class="button button1">Green</button>
	<br>
	<button class="button button2">Blue</button>
	<br>
	<a href="http://localhost:8000/user-register">Go to user register (Create new password)</a> 
	<br>
	
	</body>
	</html>`)
  
	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "amantarar01@gmail.com", "mdyrprmdvxpfxjjp")
  
	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
  
	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
	  fmt.Println(err)
	  panic(err)
	}
  
	return
}





func Create_new_password(w http.ResponseWriter, r *http.Request){

	fmt.Println("sdfmsdfm")
	claims,_:=DecodeToken(w,r,"forgotuser_cookie")

	fmt.Println("claims",claims)
	fmt.Println("helooooooooooo")



	utils.SetHeader(w)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)

		return
	}  

	input :=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	
	err := validation.Validate(input,
		validation.Map(
			
			
			validation.Key("new_password",validation.Required,validation.Length(8, 16),),

			validation.Key("confirm_new_password",validation.Required,validation.Length(8, 16),),

			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}else if input["new_password"]!=input["confirm_new_password"]{


		Res.Response("Bad Request",400,"Password not matching","",w)
		return
	}

	//change the users password to new password

	var user models.User
	fmt.Println("input email",input["email"])
	err1:=db.DB.Where("email=?",claims.Email).First(&user).Error
	if err1!=nil{
		Res.Response("Server error",500,err1.Error(),"",w)
		return

	}
	fmt.Println("user old",user)

	user.Password,_=HashPassword(input["confirm_new_password"])   

	fmt.Println("user new",user)

	fmt.Println("claims.email",claims.Email)         
	er:=db.DB.Where("email=?",claims.Email).Updates(&user).Error
	if er!=nil{
		Res.Response("Server error",500,er.Error(),"",w)
		return

	}


	Res.Response("Success",200,"Password Set Successfully","",w)



}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}






func User_login_with_password(w http.ResponseWriter, r *http.Request){

	utils.SetHeader(w)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)

		return
	}  


	input :=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	
	err := validation.Validate(input,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("email", validation.Required, is.Email),
			
			validation.Key("password",validation.Required,validation.Length(8, 16),),
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}


	var user models.User

	er:=db.DB.Where("email = ?", input["email"]).First(&user).Error

	if er!=nil{
		Res.Response("Server error ",500,"email not present","",w)

	}
	fmt.Println("password input",input["password"])
	


	
	fmt.Println("password in db",user.Password)
	
	er1 := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input["password"]))

	if er1==nil{

		//give user its cookie with token

			// jwt authentication token
			expirationTime := time.Now().Add(2 * time.Minute)
			fmt.Println("expiration time is: ", expirationTime)
	
			// check if the user is valid then only create token
			claims := models.Claims{
				Phone: user.Contact_no,
				User_id:user.User_id,
				Active:true,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(expirationTime),
				},
			}
			fmt.Println("claims: ", claims)
	
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			fmt.Println("token: ", token)

			tokenString, err := token.SignedString((con.Jwt_key))

			if err != nil {
				fmt.Println("error is :", err)
				w.WriteHeader(http.StatusInternalServerError)
			Res.Response("Bad gateway",500,err.Error(),"",w)
				}
			fmt.Println("tokenString",tokenString)
			
			
				//put this token in a cookie

				http.SetCookie(w, &http.Cookie{
					Name:    "token",
					Value:   tokenString,
					HttpOnly: true,
					Expires: expirationTime.Add(8*time.Hour),
				})
			
				fmt.Println("cookie set hua")
			
			if err != nil {
				
				
				// response.Message="Invalid or expired token"
				Res.Response("unauthorized",401,"Invalid or expired token","",w)

				
				fmt.Println("invalid token",err)
			
			}
			//after the token is provided
	
	
				user.Token=tokenString
				user.LoggedIn=claims.Active
			
			
			fmt.Println("token provided successfully")
			Res.Response("OK",200,"token provided successfully AND PHONE NUMBER VERIFIED","",w)
			


	}else {
		// w.Write([]byte("Verifictaion failed"))
		fmt.Println("verification failed")
		
		Res.Response("Unauthorized",401,er1.Error()+"Password wrong or user not present","",w)
	}

	



}




// @Description User login with name and contact_no
// @Accept json
// @Produce json
//  @Param  details body string true "name and contact number" SchemaExample({"name":"john doe","contact":"1234567890"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /user-login-contact [post]
func User_login_with_contact_no(w http.ResponseWriter,r *http.Request){

	
	utils.SetHeader(w)
	

	EnableCors(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)

		return
	}  
	
	w.Header().Set("Content-Type", "application/json")
	
	//take the user name ,email and contact number

	
	input :=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	
	err := validation.Validate(input,
		validation.Map(
			//contact number must be between 8-10 digits.
			validation.Key("contact",validation.Required,validation.Length(8, 10),is.Digit,),
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}

	

	var user models.User
	user.Contact_no=input["contact"]
	


		//IF ALREADY EXISTS JUST UPDATE THE TOKEN OTHERWISE CREATE NEW USER

		query:="select exists(select * from users where contact_no='"+input["contact"]+"');"
		var user_exists bool
		db.DB.Raw(query).Scan(&user_exists)
		if !user_exists{

			//create entry 
			er:=db.DB.Create(&user).Error
			
			if er!=nil{
			Res.Response("server error",500,er.Error(),"",w)
			return
	
			}
	
		}


	//send otp according to the contact number entered
	errr:=sendOtp("+91" + input["contact"],w)
	if errr!=nil{

		return
	}
	//generate an Otp
	Res.Response("OK",200,"OTP SENT","",w)


}












// @Description User logout
// @Produce json
// @Tags User
// @Success 200 {object} models.Response
// @Router /user-logout [get]
func User_logOut(w http.ResponseWriter, r *http.Request){


	utils.SetHeader(w)
	

	var user models.User




	
	


	
	userData := r.Context().Value("user")
	
	var userDetails models.Claims

	userDetails=userData.(models.Claims)

	user.User_id=userDetails.User_id
	


	
	c, err := r.Cookie("token")

	if err!=nil{
		Res.Response("Unauthorized",401,"Valid Cookie not found","",w)
		
	}

	var blacklist_token models.Blacklisted_tokens
	blacklist_token.Token=c.Value
	db.DB.Create(&blacklist_token)
	fmt.Println("token blacklist hua")
	user.LoggedIn=false
	db.DB.Where("user_id=?",user.User_id).Updates(&user)


	
	cookie:=&http.Cookie{

		Name: "token",
		Expires: time.Now(),
	}

	http.SetCookie(w,cookie)

	fmt.Println("expired cookie set hua")
	fmt.Fprint(w,"cookie expire hua?")

	Res.Response("Success",200,"Logged out successfully","",w)



	//overwrite with a just in time expired cookie

}



// @Description Get song by id
// @Accept json
// @Produce json
//  @Param  details body string true "enter song id" SchemaExample({"song_id":"xyz"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /get-song [post]
func GetSong(w http.ResponseWriter,r * http.Request){


	utils.SetHeader(w)


	// //get the song from db based on the id of song
	if r.Method != http.MethodPost {
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		return
	}


	w.Header().Set("Content-Type", "application/json")

	

	input :=make(map[string]string)

	json.NewDecoder(r.Body).Decode(&input)

	er := validation.Validate(input,
		validation.Map(
			// id cannot be empty
			validation.Key("id", validation.Required,),
		
		),
	)
	
	if er!=nil{

		Res.Response("Bad Request",400,er.Error(),"",w)
		return
	}

	var song models.AudioFile
	song.ID=input["id"]

	err:=db.DB.Where("id=?",input["id"]).First(&song).Error
	if err!=nil{

		Res.Response("server error",500,err.Error(),"",w)
		return

	}



	Res.Response("OK",200,"Success",song,w)


}

// @Description Get all song
// @Produce json
// @Tags User
// @Success 200 {object} models.Response
// @Router /get-AllSongs [get]
func Get_All_Songs(w http.ResponseWriter,r *http.Request){

	utils.SetHeader(w)

	var songs []models.AudioFile



	query:="SELECT * FROM audio_files;"
	err:=db.DB.Raw(query).Scan(&songs).Error
	if err!=nil{
		Res.Response("server error",500,err.Error(),"",w)
		return

	}
	Res.Response("OK",200,"Success",songs,w)

	// json.NewEncoder(w).Encode(&songs)



}

// @Description Create playlist
// @Accept json
// @Produce json
//  @Param  details body string true "enter playlist name and song id" SchemaExample({"playlist_name":"name of your playlist","song_id":"xyz"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /create-playlist [post]
func CreatePlaylist(w http.ResponseWriter,r * http.Request){


	//custom playlist
	//user want to add songs to his/her playlist

	//it will take  playlist_name and path 
	//user_id will be set from the token
	utils.SetHeader(w)


	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		return
		
	}
	




	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)

	er := validation.Validate(input,
		validation.Map(
			// playlist cannot be empty
			validation.Key("playlist_name", validation.Required),
			validation.Key("song_id",validation.Required),
		
		),
	)
	
	if er!=nil{

		Res.Response("Bad Request",400,er.Error(),"",w)
		return
	}


	//extract the user_id from the token
	
	
	userData := r.Context().Value("user")
	
	var userDetails models.Claims

	userDetails=userData.(models.Claims)

	input["user_id"]=userDetails.User_id
	


	
	var playlist models.Playlist
	playlist.Playlist_name=input["playlist_name"]
	playlist.Song_id=input["song_id"]
	playlist.User_id=input["user_id"]

	query:="select exists(select * from playlists where user_id='"+playlist.User_id+"' and song_id='"+playlist.Song_id+"');"
	var exists bool
	db.DB.Raw(query).Scan(&exists)
	if exists{
		Res.Response("Bad Request",400,"already exists","",w)
		return 
	}

	db.DB.Create(&playlist)

	Res.Response("OK",200,"added to playlist","",w)


}

// @Description Show existing playlist
// @Produce json
// @Tags User
// @Success 200 {object} models.Response
// @Router /show-playlist [get]
func Show_playlist(w http.ResponseWriter, r *http.Request){


	//take the name of the playlist
	//userid will be automatically fetch from token
	utils.SetHeader(w)


	var playlist models.Playlist //for decoding body
	var playlists_songs []models.Playlist //for display all songs of playlist


	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)

	err := validation.Validate(input,
		validation.Map(
			// playlist cannot be empty
			validation.Key("playlist_name", validation.Required),
		
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}
	playlist.Playlist_name=input["playlist_name"]
	
	userData := r.Context().Value("user")
	
	var userDetails models.Claims

	userDetails=userData.(models.Claims)

	playlist.User_id=userDetails.User_id


	

	fmt.Println("plalist.playlist name",playlist.Playlist_name)
	fmt.Println("",playlist.User_id)

	
	query:="SELECT * FROM playlists WHERE playlist_name='"+ playlist.Playlist_name +"'AND user_id='"+ playlist.User_id+"';"



	er:=db.DB.Raw(query).Scan(&playlists_songs).Error
	if er!=nil{
		Res.Response("server error",500,er.Error(),"",w)
		return
	}
	if playlists_songs==nil{

		Res.Response("Bad request",400,"Does not Exists","",w)
		return
	}
	fmt.Println("playlist songs",playlists_songs)

	Res.Response("OK",200,"Success",playlists_songs,w)
	


}

// @Description Add your Fav. Song with this api
// @Accept json
// @Produce json
//  @Param  details body string true "enter song id" SchemaExample({"id":"xyz"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /add-fav-song [post]
func Add_song_toFav(w http.ResponseWriter,r *http.Request){



	utils.SetHeader(w)

	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}

	
	//take the song id as input from r.body
	

	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)

	err1 := validation.Validate(input,
		validation.Map(
			// id cannot be empty
			
			validation.Key("id",validation.Required),
		
		),
	)
	
	if err1!=nil{

		Res.Response("Bad Request",400,err1.Error(),"",w)
		return
	}

	

	var fav_song models.Fav_Songs



	userData := r.Context().Value("user")
	
	var userDetails models.Claims

	userDetails=userData.(models.Claims)

	fav_song.User_id=userDetails.User_id



	
	fav_song.Song_id=input["id"]//song ka path set hogya

	query:="select exists(select * from fav_songs where user_id='"+fav_song.User_id+"' and song_id='"+fav_song.Song_id+"');"
	var exists bool
	db.DB.Raw(query).Scan(&exists)
	if exists{
		Res.Response("Bad Request",400,"already exists","",w)
		return 
	}
	er:=db.DB.Create(&fav_song).Error // table me create hogya
	if er != nil {

		Res.Response("server error",500,er.Error(),"",w)
		return
	}

	Res.Response("OK",200,"added to fav.","",w)

}


// @Description Add your Fav. Song with this api
// @Accept json
// @Produce json
// @Tags User
// @Success 200 {object} models.Response
// @Router /get-fav-song [get]
func Get_Fav_song_list(w http.ResponseWriter,r * http.Request){

	utils.SetHeader(w)


	var user models.User

	var fav_songs_list []models.Fav_Songs



	userData := r.Context().Value("user")
	
	var userDetails models.Claims

	userDetails=userData.(models.Claims)

	user.User_id=userDetails.User_id

	

	query:="SELECT * FROM fav_songs;"

	db.DB.Raw(query).Scan(&fav_songs_list)

	Res.Response("OK",200,"success",fav_songs_list,w)



}


// @Description Add to Recently_Played list
// @Accept json
// @Produce json
//  @Param  details body string true "enter song id" SchemaExample({"song_id":"xyz"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /add-to-recentlyPlayed [post]
func Add_to_RecentlyPlayed(w http.ResponseWriter,r *http.Request){


	//userid will be fetch from token
	//playedAt time will be set manually
	//song id will be decoded from r.body
	utils.SetHeader(w)


	if r.Method != http.MethodPost {
		
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}

	

	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)

	er := validation.Validate(input,
		validation.Map(
			// song id cannot be empty
			validation.Key("song_id",validation.Required),
		
		),
	)
	
	if er!=nil{

		Res.Response("Bad Request",400,er.Error(),"",w)
		return
	}

	

		var recent_add_song models.Recently_Played
		recent_add_song.PlayedAt=time.Now() //played at this time 
		

		


		userData := r.Context().Value("user")
	
		var userDetails models.Claims
	
		userDetails=userData.(models.Claims)
	
		recent_add_song.User_id=userDetails.User_id

	

		recent_add_song.Song_id=input["song_id"]
		query:="select exists(select * from recently_playeds where user_id='"+recent_add_song.User_id+"' and song_id='"+recent_add_song.Song_id+"');"
		var exists bool
		db.DB.Raw(query).Scan(&exists)
		if exists{
			Res.Response("Bad Request",400,"already exists","",w)
			return 
		}
		
		err1:=db.DB.Create(&recent_add_song).Error
		if er!=nil{
			Res.Response("server error",500,err1.Error(),"",w)
			return
		}
		
		Res.Response("OK",200,"added to recently played","",w)
		

	
}

// @Description Get Recently_Played list
// @Accept json
// @Produce json
// @Tags User
// @Success 200 {object} models.Response
// @Router /get-recently-playedList [get]
func Get_Recently_Played_Songs(w http.ResponseWriter,r *http.Request){

	//get the list of 20 recently played songs
	utils.SetHeader(w)


	//get the userid from token parsing-->(to know who is asking for this request)
	//filter the entries of recently played based on userid

	var user models.User


	userData := r.Context().Value("user")
	
	var userDetails models.Claims

	userDetails=userData.(models.Claims)

	user.User_id=userDetails.User_id

	

	query:="SELECT * FROM recently_playeds WHERE user_id='"+user.User_id+"' ORDER BY  played_at DESC LIMIT 20"
	var list_of_recently_played []models.Recently_Played
	er:=db.DB.Raw(query).Scan(&list_of_recently_played).Error
	if er!=nil{

		Res.Response("server error",500,er.Error(),"",w)
		return
	}

	Res.Response("OK",200,"success",list_of_recently_played,w)
	



}


// @Description Get Artist
// @Accept json
// @Produce json
//  @Param  details body string true "enter name of artist SchemaExample({"artist_name":"Arijit Singh"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /get-artist [post]
func Get_Artist(w http.ResponseWriter,r *http.Request){


	//get the name of the artist from r.body
	utils.SetHeader(w)
	
	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)

	er := validation.Validate(input,
		validation.Map(
			//artist name required
			validation.Key("artist_name",validation.Required),
		
		),
	)
	
	if er!=nil{

		Res.Response("Bad Request",400,er.Error(),"",w)
		return
	}

	

	//based on the artist name get the songs from audiofiles table
	var artist models.Artist
	artist.Artist_name=input["artist_name"]

	query:="SELECT * FROM audio_files WHERE artist='"+artist.Artist_name+"';"

	var list_of_artist_songs []models.AudioFile

	err:=db.DB.Raw(query).Scan(&list_of_artist_songs).Error
	if err!=nil{

		Res.Response("server error",500,err.Error(),"",w)
		return
	}


	Res.Response("OK",200,"Success",list_of_artist_songs,w)
}


// @Description Get Album
// @Accept json
// @Produce json
// @Param  details body string true "enter name of album SchemaExample({"album_name":"xyzBest"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /get-album [post]
func Get_Album(w http.ResponseWriter,r *http.Request){

	utils.SetHeader(w)

		//get album based on the name of album
		if r.Method != http.MethodPost {
		
			Res.Response("Method Not Allowed ",405,"use correct http method","",w)
			
		}
		

		var album models.Album
		var album_song_list []models.Album

		input:=make(map[string]interface{})
		json.NewDecoder(r.Body).Decode(&input)
	
		er := validation.Validate(input,
			validation.Map(
			
				validation.Key("album_name",validation.Required),
			
			),
		)
		
		if er!=nil{
	
			Res.Response("Bad Request",400,er.Error(),"",w)
			return
		}

		json.NewDecoder(r.Body).Decode(&album)

		query:="SELECT * FROM albums WHERE album_name='"+album.Album_name+"';"
		err:=db.DB.Raw(query).Scan(&album_song_list).Error
		if err!=nil{

			Res.Response("server error",500,err.Error(),"",w)
			return
		}
		Res.Response("OK",200,"Success",album_song_list,w)
		
	
}


// @Description Search Song by name
// @Accept json
// @Produce json
// @Param  details body string true "enter name of song SchemaExample({"name":"song_name"})
// @Tags User
// @Success 200 {object} models.Response
// @Router /search-songs [post]
func Search_Song(w http.ResponseWriter,r *http.Request){

	utils.SetHeader(w)

	//get album based on the name of album
	if r.Method != http.MethodPost {
	
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}
	var Song_list []models.AudioFile
	
	input:=make (map[string]string)

	
	json.NewDecoder(r.Body).Decode(&input)

	er := validation.Validate(input,
		validation.Map(
		
			validation.Key("name",validation.Required),
		
		),
	)
	
	if er!=nil{

		Res.Response("Bad Request",400,er.Error(),"",w)
		return
	}


	query:="SELECT * FROM audio_files WHERE LOWER(name) LIKE LOWER('"+input["name"]+"%')UNION SELECT * FROM audio_files WHERE LOWER(name) LIKE LOWER('%"+input["name"]+"%')AND LOWER(name) NOT LIKE LOWER('"+input["name"]+"%') ;"

	err:=db.DB.Raw(query).Scan(&Song_list).Error

	if err!=nil{

		Res.Response("Bad Request",400,er.Error(),"",w)
		return
	}



	Res.Response("OK",200,"Success",Song_list,w)
	

}

