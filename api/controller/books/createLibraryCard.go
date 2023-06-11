package books

import (
	"context"
	"fmt"
	"io"
	"main/helper"
	"net/http"
	"os"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

type InsertDBPayload struct {
	Name        string   `bson:"name,omitempty"`
	Age         string   `bson:"Age,omitempty"`
	Location    string   `bson:"Location,omitempty"`
	UploadPaths []string `bson:"image_paths,omitempty"`
	UserId      string   `bson:"User-ID,omitempty"`
}

func CreateLibraryCard(w http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := request.PostFormValue("name")
	age := request.PostFormValue("age")
	location := request.PostFormValue("location")
	paths := upload(w, request)

	var database = helper.ConnectDB()
	collection := database.Collection("Users")

	var document InsertDBPayload
	document.UploadPaths = paths
	document.Name = name
	document.Age = age
	document.Location = location

	count, err := collection.Distinct(context.TODO(), "User-ID", bson.D{})
	// CountDocuments(context.TODO(), bson.D{{Key: "User-ID", Value: bson.D{{Key: "$exists", Value: true}}}})
	if err != nil {
		var p = helper.Params{
			ResponseWriter: w,
			Error:          err,
			CustomMessage:  "",
		}
		helper.GetError(p)
		return
	}

	document.UserId = strconv.Itoa(len(count) + 1)

	result, err := collection.InsertOne(context.TODO(), document)

	if err != nil {
		var p = helper.Params{
			ResponseWriter: w,
			Error:          err,
			CustomMessage:  "",
		}
		helper.GetError(p)
		return
	}

	fmt.Println(result)
	fmt.Fprintf(w, "Upload successful")

	w.WriteHeader(200)
	return
}

func upload(w http.ResponseWriter, request *http.Request) []string {
	files := request.MultipartForm.File["photo"]
	var uploadPaths []string
	for _, fileHeader := range files {
		// Restrict the size of each uploaded file to 1MB.
		// To prevent the aggregate size from exceeding
		// a specified value, use the http.MaxBytesReader() method
		// before calling ParseMultipartForm()
		if fileHeader.Size > 1024*1024 {
			http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 1MB in size", fileHeader.Filename), http.StatusBadRequest)
			return []string{}
		}

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return []string{}
		}

		defer file.Close()

		buff := make([]byte, 512)
		_, err = file.Read(buff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return []string{}
		}

		filetype := http.DetectContentType(buff)
		if filetype != "image/jpeg" && filetype != "image/png" {
			http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			return []string{}
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return []string{}
		}

		err = os.MkdirAll("./uploads", os.ModePerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return []string{}
		}

		f, err := os.Create(fmt.Sprintf("./uploads/%s", fileHeader.Filename))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return []string{}
		}
		fmt.Println(uploadPaths)
		uploadPaths = append(uploadPaths, f.Name())
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return []string{}
		}
	}
	return uploadPaths
}
