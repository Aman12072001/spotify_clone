package auth

import (
	"context"
	Resp "main/Response"
	cont "main/controllers"
	"net/http"
)




func IsAuthorizedUser(endpoint func(http.ResponseWriter,*http.Request))http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

		claims,err:=cont.DecodeToken(w,r)
		
		if claims.Valid()==nil{

			userData := context.WithValue(r.Context(), "user", claims)
			// mpData := r.Context().Value("editUser")
			endpoint(w, r.WithContext(userData))

		}
		if err!=nil{

			Resp.Response("OK",200,"token provided successfully","",w)
			
		}
	
	})
}	