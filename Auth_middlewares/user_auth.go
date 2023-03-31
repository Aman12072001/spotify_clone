package auth

import (
	"fmt"
	cons "main/utils"
	"net/http"

	jwt "github.com/golang-jwt/jwt"
)


func isAuthorized(endpoint func(http.ResponseWriter,*http.Request))http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		if r.Header["Token"] != nil {

			token,err:=jwt.Parse(r.Header["Token"][0],func(token *jwt.Token) (interface{},error){

				if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil,fmt.Errorf("error")
				}
				return cons.Jwt_key,nil
			})
			if err!=nil{
				fmt.Fprintf(w,"error in parsing")
			}
			if token.Valid {
				endpoint(w,r)
			}

	}else{

		fmt.Fprint(w,"NOt authorized")
	}
	
	})
}	