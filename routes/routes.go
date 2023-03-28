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

	log.Fatal(http.ListenAndServe(":8000", mux))

}