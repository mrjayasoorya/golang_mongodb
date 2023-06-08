// https://www.mongodb.com/docs/manual/reference/operator/query/
package books

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/helper"
	"main/models"
	"net/http"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var collection = helper.ConnectDB()

type apiCallParameters struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	payload        models.GetBooksPayload
	datasetlimit   int
}

func RequestBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Request body

	decoder := json.NewDecoder(r.Body)

	// Books model
	var p models.GetBooksPayload

	err := decoder.Decode(&p)
	if err != nil {
		var p = helper.Params{
			ResponseWriter: w,
			Error:          err,
			CustomMessage:  "Bad Payload",
		}
		helper.GetError(p)
		return
	}

	// we created Book array
	var books []models.Book

	// --------->
	// REQUEST PAYLOAD PARSING
	datasetlimit := int(10)

	if len(p.Limit) > 0 {
		intVar, err := strconv.Atoi(p.Limit)
		if err != nil {
			var p = helper.Params{
				ResponseWriter: w,
				Error:          err,
				CustomMessage:  "",
			}
			helper.GetError(p)
			return
		}
		datasetlimit = intVar

	}

	params := apiCallParameters{
		responseWriter: w,
		request:        r,
		payload:        p,
		datasetlimit:   datasetlimit,
	}
	if len(p.FreeText) > 0 {
		books = apiCall(params)
	} else {
		books = apiCall(params)
	}

	var response models.GetBooksResponsePayload

	// <-------

	// ----->
	// sending user the response
	if len(books) == 0 {
		response.Results = books
		response.Count = len(books)
		response.Error.ErrorCode = 1
		response.Error.Message = "Error"
	} else {
		response.Results = books
		response.Count = len(books)
		response.Error.ErrorCode = 0
		response.Error.Message = "Success"
	}
	json.NewEncoder(w).Encode(response) // encode similar to serialize process.
	// <-----

}

func apiCall(p apiCallParameters) []models.Book {
	var collection = helper.ConnectDB()

	// we created Book array
	var books []models.Book
	// -------->
	// DATABASE QUERYING
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "Publisher", Value: -1}})
	opts.SetLimit(int64(p.datasetlimit))

	filter := bson.D{}
	if p.payload.Title != "" && p.payload.Author != "" || p.payload.Title != "" && p.payload.Author != "" && p.payload.FreeText != "" || p.payload.Title == "" && p.payload.Author == "" && p.payload.FreeText == "" {
		// fieldval := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: p.payload.Title, Options: "xi"}}}
		// filter = bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "Book-Title", Value: fieldval}}, bson.D{{Key: "Book-Author", Value: fieldval}}}}}
		return books
	} else if p.payload.FreeText != "" {
		fieldval := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: p.payload.FreeText, Options: "xi"}}}
		filter = bson.D{{Key: "$or", Value: bson.A{bson.D{{Key: "Book-Title", Value: fieldval}}, bson.D{{Key: "Book-Author", Value: fieldval}}}}}
	} else if p.payload.Title != "" {
		filter = bson.D{{Key: "$text", Value: bson.D{{Key: "$search", Value: p.payload.Title}}}}
	} else if p.payload.Author != "" {
		fieldval := bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: p.payload.Author, Options: "xi"}}}
		filter = bson.D{{Key: "Book-Author", Value: fieldval}}
	}
	// filter = bson.D{{Key: "Book-Title", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: p.payload.Title, Options: "xi"}}}}}
	cur, err := collection.Collection("Books").Find(context.TODO(), filter, opts)
	fmt.Println(filter)

	if err != nil {
		var p = helper.Params{
			ResponseWriter: p.responseWriter,
			Error:          err,
		}
		helper.GetError(p)
		return books
	}

	defer cur.Close(context.TODO())
	// <---------

	// -------->
	// POST DATABASE QUERIYING MANUPLATING THE DATA TO JSON

	for cur.Next(context.TODO()) {

		var book models.Book
		err := cur.Decode(&book) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		books = append(books, book)
	}

	// <---------

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return books
}
