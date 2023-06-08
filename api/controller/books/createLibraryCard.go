package books

import (
	"fmt"
	"net/http"
)

func CreateLibraryCard(w http.ResponseWriter, request *http.Request) {
	err := request.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Access the name key - First Approach
	fmt.Println(request.Form["name"])
	//Access the name key - Second Approach
	fmt.Println(request.PostForm["name"])
	//Access the name key - Third Approach
	fmt.Println(request.MultipartForm.Value["name"])
	//Access the name key - Fourth Approach
	fmt.Println(request.PostFormValue("name"))
	//Access the age key - First Approach
	fmt.Println(request.Form["age"])
	//Access the age key - Second Approach
	fmt.Println(request.PostForm["age"])
	//Access the age key - Third Approach
	fmt.Println(request.MultipartForm.Value["age"])
	//Access the age key - Fourth Approach
	fmt.Println(request.PostFormValue("age"))
	//Access the photo key - First Approach
	// _, h, err := request.FormFile("photo")
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// saveFile(h, "formfile")
	// //Access the photo key - Second Approach
	for _, h := range request.MultipartForm.File["photo"] {
		// err := saveFile(h, "mapaccess")
		// if err != nil {
		// 	w.WriteHeader(http.StatusBadRequest)
		// 	return
		// }
		fmt.Println(h.Filename)
	}
	w.WriteHeader(200)
	return
}
