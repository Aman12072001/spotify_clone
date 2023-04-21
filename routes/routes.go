package routes

import (
	"fmt"
	"log"
	auth "main/Auth_middlewares"
	cont "main/controllers"
	"main/db"
	_ "main/docs"
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



	//ADMIN
	mux.Handle("/add-song",auth.IsAuthorizedAdmin(cont.Add_Song))
	mux.Handle("/add-img",auth.IsAuthorizedAdmin(cont.Add_Thumbnail_Img))
	mux.Handle("/create-album",auth.IsAuthorizedAdmin(cont.Create_Album))



	//USER AUTHENTICATION
	mux.HandleFunc("/user-login-contact",cont.User_login_with_contact_no)
	mux.HandleFunc("/verify-otp",cont.VerifyOtp)
	mux.Handle("/user-logout",auth.IsAuthorizedUser(cont.User_logOut))
	mux.HandleFunc("/user-register",cont.User_SignUp)
	mux.HandleFunc("/user-login-email",cont.User_login_with_password)
	mux.HandleFunc("/forgot-password",cont.Forgot_Password)
	mux.HandleFunc("/create-new-password",(cont.Create_new_password))


	//create playlist
	mux.Handle("/create-playlist",auth.IsAuthorizedUser(cont.CreatePlaylist))
	mux.Handle("/show-playlist",auth.IsAuthorizedUser(cont.Show_playlist))


	//Songs
	mux.HandleFunc("/get-song",cont.GetSong)
	mux.HandleFunc("/get-allSongs",cont.Get_All_Songs)
	mux.Handle("/add-fav-song",auth.IsAuthorizedUser(cont.Add_song_toFav))
	mux.Handle("/add-to-recentlyPlayed",auth.IsAuthorizedUser(cont.Add_to_RecentlyPlayed))
	mux.Handle("/get-recently-playedList",auth.IsAuthorizedUser(cont.Get_Recently_Played_Songs))
	mux.HandleFunc("/get-album",cont.Get_Album)
	mux.HandleFunc("/search-songs",cont.Search_Song)

	//Get artist
	mux.HandleFunc("/get-artist",cont.Get_Artist)
	
	//create payment
	mux.Handle("/make-payment",auth.IsAuthorizedUser(cont.MakepaymentHandler))
	mux.HandleFunc("/razorpay-response",cont.Razorpay_Response)
	

	//Swagger handler
	mux.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)



	log.Fatal(http.ListenAndServe(":8000", mux))

}