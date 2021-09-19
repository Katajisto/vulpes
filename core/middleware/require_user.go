package middleware

import (
	"log"
	"net/http"

	"vulpes.ktj.st/models"
)

type RequireUser struct {
	us *models.UserService
}

func NewRequreUserMw(us *models.UserService) *RequireUser {
	return &RequireUser{us: us}
}

func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := mw.us.BySession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		log.Println("User found: ", user)
		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.Handler {
	return mw.ApplyFn(next.ServeHTTP)
}
