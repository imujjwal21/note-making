package httptransport

import (
	"html/template"
	"net/http"
	"notemaking/notes"
	"notemaking/users"

	"github.com/gorilla/mux"
)

var templates *template.Template

func NewHandler(userStorage users.UserDataStore, noteStorage notes.NoteDataStore) http.Handler {

	router := mux.NewRouter()
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// w.Header().Set("content-type", "application/json")
			w.Header().Set("content-type", "text/html")
			next.ServeHTTP(w, r)
		})
	})

	templates = template.Must(template.ParseGlob("templates/*.html"))

	router.HandleFunc("/", indexGetHandler()).Methods("GET")

	router.HandleFunc("/register", registerGetHandler()).Methods("GET")
	router.HandleFunc("/register", registerPostHandler(userStorage)).Methods("POST")

	router.HandleFunc("/login", loginGetHandler()).Methods("GET")
	router.Path("/login").Methods(http.MethodPost).HandlerFunc(loginPostHandler(userStorage))

	router.HandleFunc("/logout", logoutGetHandler()).Methods("GET")

	//notes

	router.HandleFunc("/addnotes", addNoteGetHandler()).Methods("GET")
	router.HandleFunc("/addnotes", addNotePostHandler(noteStorage)).Methods("POST")

	router.HandleFunc("/shownotes", showNoteGetHandler(noteStorage)).Methods("GET")

	router.HandleFunc("/deletenote", deleteNoteGetHandler(noteStorage)).Methods("GET")

	router.HandleFunc("/editnote", editNoteByIdGetHandler()).Methods("GET")
	router.HandleFunc("/editnote", editNoteByIdPostHandler(noteStorage)).Methods("POST")

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	return router
}
