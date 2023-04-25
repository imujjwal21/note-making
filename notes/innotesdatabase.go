package notes

import (
	"context"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
)

type inNoteDatabase struct {
	database *mongo.Database
}

func NewInNoteDatabase(db *mongo.Database) NoteDataStore {
	return &inNoteDatabase{db}
}

const names = "Notes Package"

func (i *inNoteDatabase) CreateNote(ctx context.Context, title, content, email string) error {

	_, span := otel.Tracer(names).Start(ctx, "Create Note")
	defer span.End()

	var lastrecord bson.M
	var note Note

	collectionNote := i.database.Collection("note")

	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	err := collectionNote.FindOne(ctx, bson.M{}, opts).Decode(&lastrecord)

	if err != nil {
		note = Note{Id: 0, Title: title, Content: content, Email: email}
	} else {
		id := fmt.Sprint(lastrecord["id"])
		noteId, _ := strconv.Atoi(id)

		note = Note{Id: noteId + 1, Title: title, Content: content, Email: email}
	}

	_, err1 := collectionNote.InsertOne(ctx, note)

	if err1 != nil {
		return err1
	}
	fmt.Println("Successfully Added in DB")

	return nil
}

func (i *inNoteDatabase) GetAllNotes(ctx context.Context, email string) ([]Note, error) {

	_, span := otel.Tracer(names).Start(ctx, "Read Notes")
	defer span.End()

	var noteSlice []Note

	collectionNote := i.database.Collection("note")
	filter := bson.D{{"email", bson.D{{"$eq", email}}}}

	cursor, err := collectionNote.Find(ctx, filter)
	if err != nil {
		return noteSlice, err
	}

	if err = cursor.All(ctx, &noteSlice); err != nil {
		return noteSlice, err
	}

	return noteSlice, nil
}

func (i *inNoteDatabase) DeleteNotesById(ctx context.Context, id string) error {

	_, span := otel.Tracer(names).Start(ctx, "Delete Note")
	defer span.End()

	collectionNote := i.database.Collection("note")

	noteId, _ := strconv.Atoi(id)
	_, err := collectionNote.DeleteOne(ctx, bson.M{"id": noteId})

	if err != nil {
		return err
	}

	return nil
}

func (i *inNoteDatabase) EditNoteById(ctx context.Context, id, title, content string) error {

	_, span := otel.Tracer(names).Start(ctx, "Edit Note")
	defer span.End()

	collectionNote := i.database.Collection("note")

	noteId, _ := strconv.Atoi(id)

	filter := bson.M{"id": noteId}
	update := bson.M{"$set": bson.M{"title": title, "content": content}}

	_, err := collectionNote.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}

	return nil
}
