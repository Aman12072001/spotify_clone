package routes

import (
	"fmt"
	"log"
	auth "main/Auth_middlewares"
	cont "main/controllers"
	"main/db"
	_ "main/docs"
	"main/models"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)



func Routes(){

	
	fmt.Println("Listening on port:8000")
	mux := mux.NewRouter()

	err := db.Connect()
	if err != nil {
		panic(err)

	}

	if db.DB.Migrator().HasTable(&models.Memberships{}) {
		var membership models.Memberships
		
		query := "SELECT memberships.membership_name, memberships.price FROM memberships;"
		db.DB.Raw(query).Scan(&membership)
		if membership.Membership_name=="" {
			cont.Membership_Distribution()
		}
	}
	if db.DB.Migrator().HasTable(&models.User{}){

		//create admin
		var admin models.User
		query:="SELECT users.user_id FROM users WHERE role='admin'"
		db.DB.Raw(query).Scan(&admin)
		if admin.User_id==""{
			//create admin 
			cont.Create_Admin()
			//give token to this admin
			cont.GetToken()
		}

	}

	//ADMIN
	mux.Handle("/addSong",auth.IsAuthorized(cont.Add_Song))
	mux.Handle("/addImg",auth.IsAuthorized(cont.Add_Thumbnail_Img))
	mux.Handle("/createAlbum",auth.IsAuthorized(cont.Create_Album))



	//USER AUTHENTICATION
	mux.HandleFunc("/userLogin",cont.User_login_with_contact_no)
	mux.HandleFunc("/verifyOtp",cont.VerifyOtp)
	mux.HandleFunc("/userLogout",cont.User_logOut)

	//create playlist
	mux.HandleFunc("/createPlaylist",cont.CreatePlaylist)

	//user profile
	mux.HandleFunc("/updateProfile",cont.UpdateProfile)

	//Songs
	mux.HandleFunc("/getSong",cont.GetSong)
	mux.HandleFunc("/showPlaylist",cont.Show_playlist)
	mux.HandleFunc("/getAllSongs",cont.Get_All_Songs)
	mux.HandleFunc("/addFavSong",cont.Add_song_toFav)
	mux.HandleFunc("/addToRecentlyPlayed",cont.Add_to_RecentlyPlayed)
	mux.HandleFunc("/getRecentlyPlayedList",cont.Get_Recently_Played_Songs)
	mux.HandleFunc("/getAlbum",cont.Get_Album)

	mux.HandleFunc("/searchSongs",cont.Search_Song)

	//Get artist
	mux.HandleFunc("/getArtist",cont.Get_Artist)
	
	//create payment
	mux.HandleFunc("/makePayment",cont.MakepaymentHandler)
	mux.HandleFunc("/razorpayResponse",cont.Razorpay_Response)
	

	mux.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)



	log.Fatal(http.ListenAndServe(":8000", mux))

}