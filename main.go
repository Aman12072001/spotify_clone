package main

import (
	"fmt"
	"main/routes"

	"github.com/joho/godotenv"
	server "main/server"
)


func main(){


    fmt.Println("welcome to spotify")
    envErr:=godotenv.Load(".env")
	if envErr!=nil {

		fmt.Println("could not load environment")
	}
    server.Server()
    routes.Routes()
	
 
}