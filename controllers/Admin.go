// all the handlers that reflects the functionalities that an admin can perform
package cont

import (
	"encoding/json"
	"fmt"
	"log"
	"main/db"
	"main/models"
	"net/http"
	"os"

	"github.com/bogem/id3v2"
)


func Add_Song(w http.ResponseWriter,r *http.Request) {


	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var pathh models.Paths
	json.NewDecoder(r.Body).Decode(&pathh)

	fmt.Println("path in post req",pathh)

	 //Open the audio file
	
	 file, err := os.Open(pathh.Path)
	 if err != nil {
		fmt.Println("err in file opening ")
		 log.Fatal(err)
		 
	 }
	 defer file.Close()
 
	 tag, err := id3v2.ParseReader(file,id3v2.Options{Parse: true})
	 if err != nil {
		 log.Fatal(err)
	 }
	// Create a new AudioFile object
	var audiofile models.AudioFile

	audiofile.Path=pathh.Path

	audiofile.Name=tag.Title()
	
	
	//calculate the size of the audiofile
	fileinfo,err:=file.Stat()
	
	audiofileinBytes:=fileinfo.Size()

	audiofile.Size=float64(audiofileinBytes/(1024*1024))
	


 
	
 
	 // Create a new record in the database
	 result := db.DB.Create(&audiofile)
	 if result.Error != nil {
		 fmt.Println(result.Error)
		 return
	 }
 fmt.Fprint(w,"Audio file saved to database")
	 fmt.Println("Audio file saved to database")
 

}