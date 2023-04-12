package main

import (
	"fmt"
	cont "main/controllers"
	"main/routes"
	"os"

	// "main/controllers"
	"github.com/joho/godotenv"
)

// @title Spotify API
// @version 1.0.0
// @description API Documentation for Spotify
// @host localhost:8000
func main(){


    fmt.Println("welcome to spotify")
    envErr:=godotenv.Load(".env")
	if envErr!=nil {

		fmt.Println("could not load environment")
	}

	cont.TwilioInit(os.Getenv("TWILIO_AUTH_TOKEN"))
	// fmt.Println("",os.Getenv("TWILIO_ACCOUNT_SID"))
	// fmt.Println("",os.Getenv("TWILIO_AUTH_TOKEN"))
	
 
    routes.Routes()
	
 
}