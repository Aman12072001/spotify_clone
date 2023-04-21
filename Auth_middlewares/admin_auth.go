package auth

import (
	Resp "main/Response"
	cont "main/controllers"
	"net/http"
)


func IsAuthorizedAdmin(endpoint func(http.ResponseWriter,*http.Request))http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){


		claims,err:=cont.DecodeToken(w,r,"token")

		if claims.Valid()==nil{

			
			// mpData := r.Context().Value("editUser")
			endpoint(w, r)

		}
		if err!=nil{

			Resp.Response("Unauthorized",401,"Not a valid admin","",w)
			
		}
	

	
	
	})
}	