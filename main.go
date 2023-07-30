package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	SQL "uploader/SQL"
	"uploader/common"

	"github.com/gorilla/mux"
)

type myBody struct {
	name        string
	description string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	//r.HandleFunc("/create_login", GetAllData).Methods("POST")
	//r.HandleFunc("/login", GetAllData).Methods("POST")
	r.HandleFunc("/video", GetAllData).Methods("GET")
	r.HandleFunc("/video/{id}", GetOneData).Methods("GET")
	//r.HandleFunc("/video/{id}", GetOneData).Methods("PUT")
	r.HandleFunc("/video/{id}/file", GetVideoOneData).Methods("GET")
	r.HandleFunc("/upload_video", uploadFile).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()

}

func Home(w http.ResponseWriter, r *http.Request) {

	//fonction := "[Home]"
	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte(fmt.Sprintf("Hello word")))

}

func GetAllData(w http.ResponseWriter, r *http.Request) {

	fonction := "[GetAllData]"

	if r.Method == "GET" {
		db, err := SQL.ConnProjectUploader()
		if err != nil {
			fmt.Println(fonction+" - line 16 : error on function SQL.ConnProjectUploader , ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mapVideo := make(map[int]common.Upload)
		mapVideo, err = SQL.GetVideo(db, mapVideo)
		if err != nil {
			fmt.Println(fonction+" - line 22 : error on function SQL.GetVideo , ", err)
			return
		}

		response, err := json.MarshalIndent(mapVideo, "", "\t")
		if err != nil {
			fmt.Println(fonction+" - line 22 : error on convertiting data to json, ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
	}

}

func GetOneData(w http.ResponseWriter, r *http.Request) {

	fonction := "[GetOneData]"

	if r.Method == "GET" {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			fmt.Println(fonction+" - line 63 : error on convert string to int , ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db, err := SQL.ConnProjectUploader()
		if err != nil {
			fmt.Println(fonction+" - line 16 : error on function SQL.ConnProjectUploader , ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mapVideo := make(map[int]common.Upload)
		mapVideo, err = SQL.GetOneVideo(db, mapVideo, id)
		if err != nil {
			fmt.Println(fonction+" - line 22 : error on function SQL.GetOneVideo , ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response, err := json.MarshalIndent(mapVideo, "", "\t")
		if err != nil {
			fmt.Println(fonction+" - line 22 : error on convertiting data to json, ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)

	}

}

func GetVideoOneData(w http.ResponseWriter, r *http.Request) {

	fonction := "[GetOneData]"

	if r.Method == "GET" {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			fmt.Println(fonction+" - line 63 : error on convert string to int , ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		db, err := SQL.ConnProjectUploader()
		if err != nil {
			fmt.Println(fonction+" - line 16 : error on function SQL.ConnProjectUploader , ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mapVideo := make(map[int]common.Upload)
		mapVideo, err = SQL.GetOneVideo(db, mapVideo, id)
		if err != nil {
			fmt.Println(fonction+" - line 22 : error on function SQL.GetOneVideo , ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		baseName := path.Base(mapVideo[id].Object_video.Path)
		file, err := os.Open(mapVideo[id].Object_video.Path)
		if err != nil {
			return
		}
		http.ServeContent(w, r, baseName, time.Unix(0, 0), file)

	}

}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		fmt.Println("File Upload Endpoint Hit")
		//l'ordre est important et il est important de cloner les response pour les utiliser plusieurs fois
		r1 := r
		r2 := r
		r3 := r
		// Parse our multipart form, 200 << 20 specifies a maximum
		// upload of 200 MB files.
		r1.ParseMultipartForm(200 << 20)
		r2.ParseForm()
		r3.ParseForm()
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file

		name := r2.Form.Get("name")

		description := r3.Form.Get("description")

		bodysrtuct := myBody{
			name:        name,
			description: description,
		}

		fmt.Println(bodysrtuct.name)

		file, handler, err := r1.FormFile("myFile")

		// // err := json.Unmarshal()

		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		// Create a new file in the uploads directory
		f, err := os.OpenFile("./video/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Copy the contents of the file to the new file
		_, err = io.Copy(f, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
