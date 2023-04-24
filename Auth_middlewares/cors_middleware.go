package auth

import "net/http"


func CORSMiddleware(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Co")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if r.Method == "OPTIONS" {
			return
		}
		originalHandler.ServeHTTP(w, r)
	})
}