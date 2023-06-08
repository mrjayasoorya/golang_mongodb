package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Create Struct
type Book struct {
	Id                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ISBN              string             `json:"isbn,omitempty" bson:"isbn,omitempty"`
	BookTitle         string             `json:"title" bson:"Book-Title,omitempty"`
	Author            string             `json:"author" bson:"Book-Author,omitempty"`
	YearOfPublication string             `json:"YearOfPublication" bson:"Year-Of-Publication,omitempty"`
	Publisher         string             `json:"Publisher" bson:"Publisher,omitempty"`
}

type GetBooksResponsePayload struct {
	Results []Book   `json:"results"`
	Error   ErrorMsg `json:"error,omitempty"`
	Count   int      `json:"count,omitempty"`
}
type GetBooksPayload struct {
	Limit    string
	Title    string
	FreeText string
	Author   string
	// Age  int
}

type ErrorMsg struct {
	ErrorCode int    `json:"statusCode,omitempty"`
	Message   string `json:"message,omitempty"`
}
