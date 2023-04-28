package httptransport

import (
	"log"
	"net/http"
	jwttoken "notemaking/jwttoken"
	"notemaking/users"
	"time"
)

func indexGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		useremail := jwttoken.CheckToken(r)
		if useremail == "" {
			http.Redirect(w, r, "/login", 302)
			return
		}

		templates.ExecuteTemplate(w, "index.html", useremail)
	}
}

func registerGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		templates.ExecuteTemplate(w, "register.html", nil)
	}
}

func registerPostHandler(storage users.UserDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		email := r.PostForm.Get("useremail")
		password := r.PostForm.Get("userpassword")

		err := storage.Register(r.Context(), email, password)

		if err != nil {
			log.Printf("cannot able Created Account : %v", err)
			// w.WriteHeader(http.StatusInternalServerError)
			http.Redirect(w, r, "/register", 302)
			return
		}

		http.Redirect(w, r, "/login", 302)
	}
}

func loginGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		templates.ExecuteTemplate(w, "login.html", nil)
	}
}

func loginPostHandler(storage users.UserDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		email := r.PostForm.Get("useremail")
		password := r.PostForm.Get("userpassword")

		err := storage.Login(r.Context(), email, password)

		if err != nil {
			log.Printf("id or pass may be incorrect : %v", err)
			http.Redirect(w, r, "/login", 302)
			return
		}

		tokenString, err := jwttoken.CreateToken(email)

		if err != nil {
			log.Printf("Token Creation Error : %v \n", err)
			return
		}

		expirationTime := time.Now().Add(time.Minute * 5)
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		http.Redirect(w, r, "/", 302)
	}
}

func logoutGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Expires: time.Now(),
		})
		http.Redirect(w, r, "/login", 302)
	}
}
