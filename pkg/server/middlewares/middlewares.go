package middlewares

import (
	"log"
	"net/http"
	"time"
)

// CORS (Cross-Origin Resource Exchange)
//
// add new HTTP headers and allow requests from different domains
func CORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	}
}

// Recovery middleware that recovers the server in case of panic
func Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[Middleware] %s panic recovered:\n%s\n",
					time.Now().Format("2006/01/02 - 15:04:05"), err)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}
}

// Logger shows every request made to the server
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(`[Logger: %s] Host: %s, Uri: %s, Method: %s, Path: %s, User-Agent: %s`,
			time.Now().Format("2006/01/02 - 15:04:05"),
			r.Host,
			r.RequestURI,
			r.Method,
			r.URL.Path,
			r.UserAgent(),
		)

		next.ServeHTTP(w, r)
	}
}
