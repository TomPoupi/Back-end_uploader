package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	common "uploader/common"
	"uploader/config"
	logger "uploader/logger"
	video "uploader/video"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type myBody struct {
	name        string
	description string
}

func main() {

	Function := "[main]"
	var line int

	// mise en place des fichiers de log
	logger.Init()
	lognew, _ := logger.LogSwitch("main.log")

	//mise en place config
	config.Init(lognew)

	line = common.GetLine()
	lognew.WithFields(log.Fields{
		"Function": Function,
		"comment":  "L" + strconv.Itoa(line) + " - initialisation serveur",
	}).Info()

	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	//r.HandleFunc("/create_login", GetAllData).Methods("POST")
	//r.HandleFunc("/login", GetAllData).Methods("POST")
	r.HandleFunc("/video", video.GetAllData).Methods("GET")
	r.HandleFunc("/video/{id}", video.OperationOneData).Methods("GET")
	r.HandleFunc("/video/{id}", video.OperationOneData).Methods("PUT")
	r.HandleFunc("/video/{id}", video.OperationOneData).Methods("DELETE")
	r.HandleFunc("/video/{id}/file", video.GetVideoOneData).Methods("GET")
	r.HandleFunc("/upload_video", video.UploadVideo).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:80",
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
