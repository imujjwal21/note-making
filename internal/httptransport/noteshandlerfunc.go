package httptransport

import (
	"log"
	"net/http"
	jwttoken "notemaking/jwtToken"
	"notemaking/notes"
)

func addNoteGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		useremail := jwttoken.CheckToken(r)
		if useremail == "" {
			http.Redirect(w, r, "/login", 302)
			return
		}
		templates.ExecuteTemplate(w, "addnotes.html", useremail)
	}
}

func addNotePostHandler(storage notes.NoteDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		email := r.PostForm.Get("useremail")
		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")

		err := storage.CreateNote(r.Context(), title, content, email)

		if err != nil {
			log.Printf("Notes Adding Error : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", 302)
	}
}

func showNoteGetHandler(storage notes.NoteDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		useremail := jwttoken.CheckToken(r)
		if useremail == "" {
			http.Redirect(w, r, "/login", 302)
			return
		}

		noteSlice, err := storage.GetAllNotes(r.Context(), useremail)
		if err != nil {
			log.Printf("cannot able Created Account : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		templates.ExecuteTemplate(w, "shownotes.html", noteSlice)
	}
}

func deleteNoteGetHandler(storage notes.NoteDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		useremail := jwttoken.CheckToken(r)
		if useremail == "" {
			http.Redirect(w, r, "/login", 302)
			return
		}

		r.ParseForm()
		noteId := r.URL.Query().Get("deleteNoteId")

		err1 := storage.DeleteNotesById(r.Context(), noteId)

		if err1 != nil {
			log.Printf("cannot delete note : %v", err1)
			return
		}

		http.Redirect(w, r, "/shownotes", 302)
	}
}

func editNoteByIdGetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		useremail := jwttoken.CheckToken(r)
		if useremail == "" {
			http.Redirect(w, r, "/login", 302)
			return
		}

		r.ParseForm()
		noteId := r.URL.Query().Get("editNoteId")

		templates.ExecuteTemplate(w, "editnotes.html", noteId)
	}
}

func editNoteByIdPostHandler(storage notes.NoteDataStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()
		noteId := r.PostForm.Get("NoteId")
		title := r.PostForm.Get("title")
		content := r.PostForm.Get("content")

		err := storage.EditNoteById(r.Context(), noteId, title, content)

		if err != nil {
			log.Printf("Notes Adding Error : %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/shownotes", 302)
	}
}
