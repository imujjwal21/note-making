package users

import "context"

type User struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type UserDataStore interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) error
	CheckEmailExits(ctx context.Context, email string) error
}
