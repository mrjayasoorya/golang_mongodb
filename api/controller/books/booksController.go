package books

import (
	"context"
	"encoding/json"
	"log"
	"main/helper"
	"main/models"
	"net/http"

	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	var collection = helper.ConnectDB()
	w.Header().Set("Content-Type", "application/json")

	// Request body

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	decoder := json.NewDecoder(r.Body)
	var p models.GetBooksPayload
	err := decoder.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Do something with the Person struct...

	// we created Book array
	var books []models.Book

	// --------->
	// REQUEST PAYLOAD PARSING
	// bson.M{},  we passed empty filter. So we want to get all data.

	datasetlimit := int(10)

	if len(p.Limit) > 0 {
		intVar, err := strconv.Atoi(p.Limit)
		if err != nil {
			var errorMsg models.ErrorMsg
			errorMsg.ErrorCode = 400
			errorMsg.Message = "something went wrong!"
			json.NewEncoder(w).Encode(errorMsg)
			return
		}
		datasetlimit = intVar

	}
	// <-------

	// -------->
	// DATABASE QUERYING
	opts := options.Find().SetLimit(int64(datasetlimit))
	cur, err := collection.Collection("Books").Find(context.TODO(), bson.M{}, opts)

	if err != nil {
		var p = helper.Params{
			ResponseWriter: w,
			Error:          err,
			// CustomMessage:  "",
		}
		helper.GetError(p)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())
	// <---------

	// -------->
	// POST DATABASE QUERIYING MANUPLATING THE DATA TO JSON

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var book models.Book
		// & character returns the memory address of the following variable.
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

	// ----->
	// sending user the response
	json.NewEncoder(w).Encode(books) // encode similar to serialize process.
	// <-----
}
