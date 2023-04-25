package users

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel"
)

type inUserDatabase struct {
	database *mongo.Database
}

func NewInUserDatabase(db *mongo.Database) UserDataStore {
	return &inUserDatabase{db}
}

const names = "User Package"

func (i *inUserDatabase) Register(ctx context.Context, email, password string) error {

	_, span := otel.Tracer(names).Start(ctx, "Create User")
	defer span.End()

	userObject := User{Email: email, Password: password}
	collectionUser := i.database.Collection("user")

	_, err := collectionUser.InsertOne(ctx, userObject)

	if err != nil {
		return err
	}
	return nil
}

func (i *inUserDatabase) Login(ctx context.Context, email, password string) error {

	_, span := otel.Tracer(names).Start(ctx, "Check User")
	defer span.End()

	var userSlice []User
	collectionUser := i.database.Collection("user")

	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"email", bson.D{{"$eq", email}}}},
				bson.D{{"password", bson.D{{"$eq", password}}}},
			},
		},
	}

	cursor, err := collectionUser.Find(ctx, filter)
	if err != nil {
		return err
	}

	if err = cursor.All(ctx, &userSlice); err != nil {
		return err
	}
	if len(userSlice) == 0 {
		return errors.New("Wrong Email or Password")
	}

	return nil
}
