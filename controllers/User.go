package cont

import (
	"encoding/json"
	"fmt"
	"main/db"
	"main/models"
	"net/http"
	"os"

	"github.com/twilio/twilio-go"
	api "github.com/twilio/twilio-go/rest/api/v2010"
)


func User_login(w http.ResponseWriter,r *http.Request){

	
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	EnableCors(&w)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	json.NewDecoder(r.Body).Decode(&user)


	// Find your Account SID and Auth Token at twilio.com/console
	// and set the environment variables. See http://twil.io/secure
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")


	var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	

	params := &api.CreateMessageParams{}
	params.SetBody("This is the ship that made the Kessel Run in fourteen parsecs?")
	// params.SetFrom("+15017122661")
	params.SetTo(user.Contact_no)
	

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		if resp.Sid != nil {
			fmt.Println(*resp.Sid)
		} else {
			fmt.Println(resp.Sid)
		}
	}



}

func GetSong(w http.ResponseWriter,r * http.Request){


	// //get the song from db based on the name of song
	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusMethodNotAllowed)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")

	var song models.AudioFile

	json.NewDecoder(r.Body).Decode(&song)

	db.DB.Where("name=?",song.Name).First(&song)

	//set the priviledge for this file ,so that frontend can access it
	// set file permissions to read and write for owner, read only for group and others
		err := os.Chmod(song.Path, 0750)
		if err != nil {
			panic(err)
		}

		//return the path for the frontend dev

	song.Path="192.180.7.255"+song.Path	

	fmt.Fprint(w,"hello bhaiya")

	json.NewEncoder(w).Encode(&song)




}