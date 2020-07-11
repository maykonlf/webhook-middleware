package entities

import "github.com/maykonlf/go-devkit/pkg/types/uuid"

type Route struct {
	ID      uuid.UUID `bson:"_id"`
	URI     string
	Body    string
	Headers map[string]string
}
