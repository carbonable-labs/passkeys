package domain

import "github.com/oklog/ulid/v2"

type AuthUser struct {
	Id    string
	Email string
}

func NewUser(email string) (*AuthUser, error) {
	id := ulid.Make()

	return &AuthUser{Id: id.String(), Email: email}, nil
}
