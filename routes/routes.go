package routes

import (
	"fmt"
	"log"
	auth "main/Auth_middlewares"
	cont "main/controllers"
	"main/db"
	"main/models"
	"net/http"
)



func Routes(){

	
	fmt.Println("Listening on port:8000")
	mux := http.NewServeMux()

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



	//USER AUTHENTICATION
	mux.HandleFunc("/verify_otp",cont.VerifyOtp)

	//create playlist
	mux.HandleFunc("/create_playlist",cont.CreatePlaylist)
	mux.HandleFunc("/user_login",cont.User_login)
	mux.HandleFunc("/update_profile",cont.UpdateProfile)
	mux.HandleFunc("/get_song",cont.GetSong)
	mux.HandleFunc("/show_playlist",cont.Show_playlist)
	mux.HandleFunc("/get_all_songs",cont.Get_All_Songs)

	//create payment
	mux.HandleFunc("/create_payment",cont.Create_Payment)
	mux.HandleFunc("/verify_payment",cont.Verify_Payment)

	


	log.Fatal(http.ListenAndServe(":8000", mux))

}