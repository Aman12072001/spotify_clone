// all the handlers that reflects the functionalities that an admin can perform
package cont

import (
	"encoding/json"
	"fmt"
	Res "main/Response"
	"main/db"
	"main/models"
	con "main/utils"
	"net/http"
	"os"
	"time"

	"github.com/bogem/id3v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v4"
)



func Create_Admin(){


	var admin models.User

	admin.Role="admin"
	admin.Name="aman-admin"
	er:=db.DB.Create(&admin).Error
	if er != nil {
		fmt.Println("db error during admin creation")
	}
	
}
func GetToken(){

	// jwt authentication token
	
	expirationTime := time.Now().Add(365* 24 * time.Hour)
	fmt.Println("expiration time is: ", expirationTime)

	// check if the user is valid then only create token

	var user models.User
	er:=db.DB.Where("role=?", "admin").First(&user).Error
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
	// fmt.Println("claims: ", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token: ", token)
	tokenString, err := token.SignedString((con.Jwt_key))
	if err != nil {
		fmt.Println("error is :", err)
		// w.WriteHeader(http.StatusInternalServerError)
		
	}
	// fmt.Println("tokenString",tokenString)
	user.Token=tokenString
	user.LoggedIn=true
	Er:=db.DB.Where("role=?", "admin").Updates(&user).Error
	if Er!=nil{

		fmt.Println("db error",Er)
	}

	
}


// @Description Add Song into app
// @Accept json
// @Produce json
// @Param  details body string true "enter PATH of song SchemaExample({"path":"/home/chicmic/Downloads/"})
// @Tags Admin
// @Success 200 {object} models.Response
// @Router /addSong [post]
func Add_Song(w http.ResponseWriter,r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	//takes audiofile path from r.body

	var pathh models.Path

	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	err := validation.Validate(input,
		validation.Map(
			// song_path cannot be empty
			
			validation.Key("path",validation.Required),
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}
	pathh.Path=input["path"]

	

	 //Open the audio file
	
	 file, err := os.Open(pathh.Path)
	 if err != nil {
		fmt.Println("err in file opening ")
		//  log.Fatal(err)
		Res.Response("Bad Request",400,"Provide proper audio file path","",w)
		return
		 
	 }
	 defer file.Close()
 
	 tag, err := id3v2.ParseReader(file,id3v2.Options{Parse: true})
	 if err != nil {
		//  log.Fatal(err)
		fmt.Println("err",err)
		 Res.Response("Server error",500,"error in audio_file parsing ","",w)
		 return

	 }
	// Create a new AudioFile object
	var audiofile models.AudioFile

	audiofile.Path=pathh.Path

	audiofile.Name=tag.Title()
	audiofile.Artist=tag.Artist()
	
	
	//calculate the size of the audiofile
	fileinfo,err:=file.Stat()
	
	audiofileinBytes:=fileinfo.Size()

	audiofile.Size=float64(audiofileinBytes/(1024*1024))

	
	

 
	 // Create a new record in the database
	 er:=db.DB.Create(&audiofile).Error
	 if er != nil {
		 fmt.Println(er.Error())
		 Res.Response("Bad Request",400,er.Error(),"",w)
		 return
	 }
	// fmt.Fprint(w,"Audio file saved to database")

	 Res.Response("Success",200,"Audio file saved to database","",w)
	fmt.Println("Audio file saved to database")
 

}


// @Description Add Thumbnail for Song 
// @Accept json
// @Produce json
// @Param  details body string true "enter Song id and path of thumbnail of song SchemaExample({"id":"xyz","img_path":"/"})
// @Tags Admin
// @Success 200 {object} models.Response
// @Router /addImg [post]
func Add_Thumbnail_Img(w http.ResponseWriter,r * http.Request){


	//take input audio_file id (in which you want to add IMg)
	//and img path
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}
	

	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	err := validation.Validate(input,
		validation.Map(
			// id and img_path cannot be empty
			validation.Key("id", validation.Required),
			validation.Key("img_path",validation.Required),
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}
	
	var song models.AudioFile
	

	

	song.Img_Path=input["img_path"]
	song.ID=input["id"]
	query:="select exists(select * from audio_files where id='"+song.ID+"' and img_path='"+song.Img_Path+"');"
	var exists bool
	db.DB.Raw(query).Scan(&exists)
	if exists{
		Res.Response("Bad Request",400,"already exists","",w)
		return 
	}	
	er:=db.DB.Where("id=?",song.ID).Updates(&song).Error
	if er!=nil{

		Res.Response("server error",500,er.Error(),"",w)
	}

	
	Res.Response("OK",200,"Thumbnail added successfully","",w)

}


// @Description Create Album
// @Accept json
// @Produce json
// @Param  details body string true "enter Song id and album name SchemaExample({"song_id":"xyz","album_name":"name"})
// @Tags Admin
// @Success 200 {object} models.Response
// @Router /createAlbum [post]
func Create_Album(w http.ResponseWriter,r * http.Request){


	//take the input song_id,album_name
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		// w.WriteHeader(http.StatusMethodNotAllowed)
		Res.Response("Method Not Allowed ",405,"use correct http method","",w)
		
	}

	

	input:=make(map[string]string)
	json.NewDecoder(r.Body).Decode(&input)
	err := validation.Validate(input,
		validation.Map(
			
			validation.Key("album_name", validation.Required),
			validation.Key("song_id",validation.Required),
			
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}

	var album models.Album
	album.Album_name=input["album_name"]
	album.Song_id=input["song_id"]

	query:="select exists(select * from albums where song_id='"+album.Song_id+"' and album_name='"+album.Album_name+"');"
	var exists bool
	db.DB.Raw(query).Scan(&exists)
	if exists{
		Res.Response("Bad Request",400,"already exists","",w)
		return 
	}

	er:=db.DB.Create(&album).Error
	if er!=nil{

		Res.Response("server error",500,er.Error(),"",w)
		
	}

	// fmt.Fprint(w,"Album created")
	Res.Response("OK",200,"Album created","",w)



}

