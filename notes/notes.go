package notes

import (
	"context"
)

type Note struct {
	Id      int    `json:"Id"`
	Title   string `json:"Name"`
	Content string `json:"Content`
	Email   string `json:"Email"`
}

type NoteDataStore interface {
	CreateNote(ctx context.Context, title, content, email string) error
	GetAllNotes(ctx context.Context, email string) ([]Note, error)
	DeleteNotesById(ctx context.Context, id string) error
	EditNoteById(ctx context.Context, id, title, content string) error
}
