package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	common "uploader/common"
	"uploader/config"
	logger "uploader/logger"
	login "uploader/login"
	video "uploader/video"

	"github.com/golang-jwt/jwt/v4"
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
	r.Handle("/", authMiddlerware(http.HandlerFunc(Home)))
	//r.HandleFunc("/create_login", GetAllData).Methods("POST")
	r.HandleFunc("/login", login.Signin).Methods("POST")
	r.HandleFunc("/refreshToken", login.Refresh).Methods("POST")
	r.HandleFunc("/video", video.OperationAllVideo).Methods("GET")
	r.HandleFunc("/video", video.OperationAllVideo).Methods("DELETE")
	r.HandleFunc("/video/{id}", video.OperationOneVideo).Methods("GET")
	r.HandleFunc("/video/{id}", video.OperationOneVideo).Methods("PUT")
	r.HandleFunc("/video/{id}", video.OperationOneVideo).Methods("DELETE")
	r.HandleFunc("/video/{id}/file", video.GetOneVideoFile).Methods("GET")
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
	w.Write([]byte(fmt.Sprintf("Hello %s !", r.Header.Get("User"))))

}

func authMiddlerware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// prend le token actuelle
		bearToken := r.Header.Get("Authorization")
		if bearToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// vérfication de la validité du token
		strArr := strings.Split(bearToken, " ")
		tknStr := strArr[1]
		claims := &login.Claims{}
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return login.JwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		r.Header.Set("User", claims.Username)
		next.ServeHTTP(w, r)
	})
}
