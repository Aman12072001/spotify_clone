package routes

import (
	"fmt"
	"log"
	cont "main/controllers"
	"main/db"
	"net/http"
)



func Routes(){


	fmt.Println("Listening on port:8000")
	mux := http.NewServeMux()

	err := db.Connect()
	if err != nil {
		panic(err)

	}

	mux.HandleFunc("/addSong",cont.Add_Song)
	mux.HandleFunc("/user_login",cont.User_login)
	mux.HandleFunc("/get_song",cont.GetSong)


	//USER AUTHENTICATION
	mux.HandleFunc("/send_otp",cont.SendOTP)
	mux.HandleFunc("/check_otp",cont.CheckOTP)

	log.Fatal(http.ListenAndServe(":8000", mux))

}