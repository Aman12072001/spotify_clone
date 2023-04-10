package cont

import (
	"encoding/json"
	"fmt"
	"main/db"
	"main/models"
	con "main/utils"
	"net/http"
	"os"
	"time"

	Res "main/Response"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/verify/v2"
	"golang.org/x/crypto/bcrypt"
)

var twilioClient *twilio.RestClient
	
func TwilioInit(password string)  {
	twilioClient = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username:con.TWILIO_ACCOUNT_SID,
		Password: password,
	})


}

func sendOtp(to string,w http.ResponseWriter) error {
	params := &openapi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	
	fmt.Println("to",to)
	// fmt.Println("from constant",con.TWILIO_AUTH_TOKEN)
	fmt.Println("from env",os.Getenv("TWILIO_AUTH_TOKEN"))

	resp, err := twilioClient.VerifyV2.CreateVerification(con.VERIFY_SERVICE_SID, params)

	if err != nil {
		fmt.Println("otp sent failed ,api error")
		fmt.Println(err.Error())
		Res.Response("Bad request",400,"TWILIO API error","",w)

		return err
		
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
	}
	return nil
}
func checkOtp(to string, code string) bool {
	params := &openapi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(code)
	
	resp, err := twilioClient.VerifyV2.CreateVerificationCheck(con.VERIFY_SERVICE_SID, params)

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
	
	err := validation.Validate(otp,
		validation.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			validation.Key("otp", validation.Required, validation.Length(6,6)),
			
			validation.Key("phone",validation.Required,validation.Length(10, 10),is.Digit,),
			
		),
	)
	
	if err!=nil{

		Res.Response("Bad Request",400,err.Error(),"",w)
		return
	}
	var user models.User
	db.DB.Where("contact_no=?",otp["phone"]).First(&user)
	if checkOtp("+91"+otp["phone"], otp["otp"]) {
		w.Write([]byte("Phone Number verified sucessfully"))

				// jwt authentication token
				expirationTime := time.Now().Add(3650*100 * time.Hour)
				fmt.Println("expiration time is: ", expirationTime)
		
				// check if the user is valid then only create token
				claims := models.Claims{
					Phone: user.Contact_no,
					User_id:user.User_id,
					Active:true,
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
				Res.Response("Bad gateway",500,err.Error(),"",w)

					
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
					
					
					// response.Message="Invalid or expired token"
				Res.Response("unauthorized",401,"Invalid or expired token","",w)

					
					fmt.Println("invalid token",err)
				
				}
				//if the token is valid
		
				//give the token string to the user so that user can validate its identity in future requests
				if claims, ok := parsedToken.Claims.(*models.Claims); ok && parsedToken.Valid {
					fmt.Printf("token will expire at :%v",  claims.ExpiresAt)
					if claims.Active{
					user.Token=parsedToken.Raw
					user.LoggedIn=claims.Active
					}
				} else {
					fmt.Println(err)
				}
				

				er:=db.DB.Where("contact_no=?",otp["phone"]).Updates(&user).Error
				if er!=nil{

					// response.Message=er.Error()
					// response.Code=500
				Res.Response("Bad gateway",500,er.Error(),"",w)

				}
				fmt.Println("token provided successfully")
				Res.Response("OK",200,"token provided successfully AND PHONE NUMBER VERIFIED","",w)
				

	} else {
		// w.Write([]byte("Verifictaion failed"))
		fmt.Println("verification failed")
		
		Res.Response("Unauthorized",401,"Verifictaion failed","",w)
	}




}

//email and password authentication

// func HashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func User_login_with_Email(w http.ResponseWriter, r *http.Request) {


	
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	







}



func DecodeToken(tokenString string) (models.Claims, error) {
	claims := &models.Claims{}

	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error")
		}
		return []byte(os.Getenv("Jwt_key")), nil
	})

	// fmt.Println("claims",claims)
	// fmt.Println("raw token",parsedToken.Raw)

	if err != nil || !parsedToken.Valid || Is_Blacklisted(parsedToken.Raw){

		// fmt.Println("fatt gya token parsing")
		return *claims, fmt.Errorf("Invalid or expired token")
	}



	return *claims, nil
}

func Is_Blacklisted(token string)bool{


	query:="SELECT EXISTS(SELECT * FROM blacklisted_tokens WHERE token='"+token+"')"
	var blacklisted_token bool
	db.DB.Raw(query).Scan(&blacklisted_token)
	if blacklisted_token{
		return true
	}

	return false
}

