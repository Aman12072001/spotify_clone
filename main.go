package main

import (
	"fmt"
	"main/routes"

	// "main/controllers"
	"github.com/joho/godotenv"
)


func main(){


    fmt.Println("welcome to spotify")
    envErr:=godotenv.Load(".env")
	if envErr!=nil {

		fmt.Println("could not load environment")
	}
	// fmt.Println("",os.Getenv("TWILIO_ACCOUNT_SID"))
	
 
    routes.Routes()
	
 
}