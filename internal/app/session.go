package app

import "github.com/google/uuid"

type Session struct {
	ID    string
	Title string
}

func NewSession() Session {
	id := uuid.New()

	return Session{
		ID: id.String(),
	}
}
