package cont

import (
	"encoding/json"
	"fmt"
	con "main/utils"
	"net/http"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
)

var client *twilio.RestClient = twilio.NewRestClientWithParams(twilio.ClientParams{
	Username:con.TWILIO_ACCOUNT_SID,
	Password: con.TWILIO_AUTH_TOKEN,
})

func sendOtp(to string) {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")


	resp, err := client.VerifyV2.CreateVerification(con.VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
	}
}
func checkOtp(to string, code string) bool {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	resp, err := client.VerifyV2.CreateVerificationCheck(con.VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println("Error is :", err)
		return false
	} else if *resp.Status == "approved" {
		return true
	} else {
		return false
	}

	// jwt authentication token
	expirationTime := time.Now().Add(10 * time.Minute)
	fmt.Println("expiration time is: ", expirationTime)

	// check if the user is valid then only create token
	claims := models.Claims{
		Phone: phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	fmt.Println("claims: ", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token: ", token)
	tokenString, err := token.SignedString(cons.JwtKey)
	if err != nil {
		fmt.Println("error is :", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// w.Write([]byte(fmt.Sprint("Token is:", tokenString)))
}

func SendOTP(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	phNumber := r.URL.Query().Get("number")

	sendOtp("+91" + phNumber)

	w.Write([]byte("OTP sent successfully "))

}

func CheckOTP(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	//phNumber := r.URL.Query().Get("number")
	var otp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&otp)
	if checkOtp("+91"+otp["phone"], otp["otp"]) {
		w.Write([]byte("Phone Number verified sucessfully"))

		
	} else {
		w.Write([]byte("Verifictaion failed"))
	}
}