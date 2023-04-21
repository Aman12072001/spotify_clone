package auth

// func IsAuthorizedSuperUser(endpoint func(http.ResponseWriter,*http.Request)) http.Handler{

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){

// 		claims,_:=cont.DecodeToken(w,r)

// 		if claims.Role=="superuser"{

// 			endpoint(w,r)
// 		}else{

// 			Resp.Response("UNAUTHORIZED",401,"user invalid","",w)
// 		}

// 	})

// }