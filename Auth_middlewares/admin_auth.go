package auth

import (
	"fmt"
	"main/models"
	cons "main/utils"
	"net/http"

	jwt "github.com/golang-jwt/jwt"
)


func IsAuthorized(endpoint func(http.ResponseWriter,*http.Request))http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		if r.Header["Token"] != nil {

			token,err:=jwt.ParseWithClaims(r.Header["Token"][0],&models.Claims{},func(token *jwt.Token) (interface{},error){

				if _,ok:=token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil,fmt.Errorf("error")
				}
				return cons.Jwt_key,nil
			})


			if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
				
				if claims.Role=="admin"{

					endpoint(w,r)
				}
			} else {
				fmt.Println(err)
			}

	}else{

		fmt.Fprint(w,"NOt authorized")
	}
	
	})
}	