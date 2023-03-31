package cont

import (
	"encoding/json"
	"fmt"
	"main/db"
	"main/models"
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
		fmt.Println("otp sent failed")
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


}

func VerifyOtp(w http.ResponseWriter,r * http.Request){

	w.Header().Set("Content-Type", "application/json")
	var otp = make(map[string]string)
	json.NewDecoder(r.Body).Decode(&otp)

	var user models.User
	db.DB.Where("contact_no=?",otp["phone"]).First(&user)
	if checkOtp("+91"+otp["phone"], otp["otp"]) {
		w.Write([]byte("Phone Number verified sucessfully"))

				// jwt authentication token
				expirationTime := time.Now().Add(100 * time.Hour)
				fmt.Println("expiration time is: ", expirationTime)
		
				// check if the user is valid then only create token
				claims := models.Claims{
					Phone: user.Contact_no,
					User_id:user.User_id,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(expirationTime),
					},
				}
				fmt.Println("claims: ", claims)
		
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				fmt.Println("token: ", token)
				tokenString, err := token.SignedString((con.Jwt_key))
				if err != nil {
					fmt.Println("error is :", err)
					w.WriteHeader(http.StatusInternalServerError)
					
				}
				fmt.Println("tokenString",tokenString)
				
				
		
				//token parsing and verification
				
		
				parsedToken ,err := jwt.ParseWithClaims(tokenString ,&models.Claims{}, func(token *jwt.Token) (interface{}, error) {
					
					if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil,fmt.Errorf("error")
					}
					return con.Jwt_key , nil
				})

				fmt.Println("tokenstring",tokenString)
				fmt.Println("parsedtoken signature",parsedToken.Signature)
				fmt.Println("parsedtoken raw",parsedToken.Raw)

				
				if err != nil {
					
					fmt.Fprint(w, "Invalid or expired token")
					fmt.Println("invalid token",err)
				
				}
				//if the token is valid
		
				//give the token string to the user so that user can validate its identity in future requests
				if claims, ok := parsedToken.Claims.(*models.Claims); ok && parsedToken.Valid {
					fmt.Printf("token will expire at :%v",  claims.ExpiresAt)
					user.Token=parsedToken.Raw
				} else {
					fmt.Println(err)
				}
				

				db.DB.Where("contact_no=?",otp["phone"]).Updates(&user)
				fmt.Println("token provided successfully")
	} else {
		w.Write([]byte("Verifictaion failed"))
	}




}

// func SendOTP(w http.ResponseWriter, r *http.Request){

// 	// we are giving phone number as input in params
// 	//and user id 
// 	w.Header().Set("Content-Type", "application/json")
// 	phNumber := r.URL.Query().Get("number")
// 	user_id:=r.URL.Query().Get("user_id")

// 	var user models.User
// 	db.DB.Where("user_id", user_id).First(&user)

// 	sendOtp("+91" + phNumber)

// 	w.Write([]byte("OTP sent successfully "))

// }

