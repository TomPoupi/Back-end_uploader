package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	common "uploader/common"
	"uploader/config"
	Ctrl "uploader/controler"
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

	OpenAccess := 10
	RestrictedAccess := 100

	r := mux.NewRouter()
	r.Handle("/", authMiddlerware(http.HandlerFunc(Home), OpenAccess)).Methods("GET")
	// r.HandleFunc("/create_login", GetAllData).Methods("POST")
	r.Handle("/login", authMiddlerware(http.HandlerFunc(login.Signin), OpenAccess)).Methods("POST")
	r.Handle("/refreshToken", authMiddlerware(http.HandlerFunc(login.Refresh), OpenAccess)).Methods("POST")
	r.Handle("/video", authMiddlerware(http.HandlerFunc(video.OperationAllVideo), OpenAccess)).Methods("GET")
	r.Handle("/video", authMiddlerware(http.HandlerFunc(video.OperationAllVideo), RestrictedAccess)).Methods("DELETE")
	r.Handle("/video/{id}", authMiddlerware(http.HandlerFunc(video.OperationOneVideo), OpenAccess)).Methods("GET")
	r.Handle("/video/{id}", authMiddlerware(http.HandlerFunc(video.OperationOneVideo), RestrictedAccess)).Methods("PUT")
	r.Handle("/video/{id}", authMiddlerware(http.HandlerFunc(video.OperationOneVideo), RestrictedAccess)).Methods("DELETE")
	r.Handle("/video/{id}/file", authMiddlerware(http.HandlerFunc(video.GetOneVideoFile), OpenAccess)).Methods("GET")
	r.Handle("/upload_video", authMiddlerware(http.HandlerFunc(video.UploadVideo), RestrictedAccess)).Methods("POST")

	// r.HandleFunc("/create_login", GetAllData).Methods("POST")
	// r.HandleFunc("/login", login.Signin).Methods("POST")
	// r.HandleFunc("/refreshToken", login.Refresh).Methods("POST")
	// r.HandleFunc("/video", video.OperationAllVideo).Methods("GET")
	// r.HandleFunc("/video", video.OperationAllVideo).Methods("DELETE")
	// r.HandleFunc("/video/{id}", video.OperationOneVideo).Methods("GET")
	// r.HandleFunc("/video/{id}", video.OperationOneVideo).Methods("PUT")
	// r.HandleFunc("/video/{id}", video.OperationOneVideo).Methods("DELETE")
	// r.HandleFunc("/video/{id}/file", video.GetOneVideoFile).Methods("GET")
	// r.HandleFunc("/upload_video", video.UploadVideo).Methods("POST")

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

func authMiddlerware(next http.Handler, LevelAccess int) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		Function := "[authMiddlerware]"
		var line int
		var Controler Ctrl.ControlerStruct

		//-----------------------------Init Controler-----------------------------------
		err := Controler.ControlLogAndDB(w, "authMiddlerware")
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Open Video File",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//------------------------------------------------------------------------------

		//******************************OpenAccess**************************************
		if LevelAccess == 10 {

			line = common.GetLine()
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - OpenAccess Route Used",
			}).Info()

			next.ServeHTTP(w, r)

		}
		//******************************************************************************

		//**************************-RestrictedAccess***********************************
		if LevelAccess == 100 {

			line = common.GetLine()
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - RestritedAccess Route Used",
			}).Info()

			//-------------------------------Recup Token------------------------------------
			// prend le token actuelle
			bearToken := r.Header.Get("Authorization")
			if bearToken == "" {
				line = common.GetLine() - 1
				err := errors.New("Token empty , user UnAuthorized")
				Controler.LogControl.WithFields(log.Fields{
					"Function": Function,
					"comment":  "L" + strconv.Itoa(line) + " - Error token user",
					"error":    err,
				}).Error()
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return

			}
			//------------------------------------------------------------------------------

			//---------------------vérfication de la validité du token----------------------
			strArr := strings.Split(bearToken, " ")
			tknStr := strArr[1]
			claims := &login.Claims{}
			tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
				return login.JwtKey, nil
			})

			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					line = common.GetLine() - 1
					Controler.LogControl.WithFields(log.Fields{
						"Function": Function,
						"comment":  "L" + strconv.Itoa(line) + " - Error token user",
						"error":    err,
					}).Error()
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}
				line = common.GetLine() - 1
				Controler.LogControl.WithFields(log.Fields{
					"Function": Function,
					"comment":  "L" + strconv.Itoa(line) + " - Error token user",
					"error":    err,
				}).Error()
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if !tkn.Valid {
				line = common.GetLine() - 1
				Controler.LogControl.WithFields(log.Fields{
					"Function": Function,
					"comment":  "L" + strconv.Itoa(line) + " - Error token user",
					"error":    err,
				}).Error()
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			line = common.GetLine()
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Token Valid !",
			}).Info()
			//------------------------------------------------------------------------------

			if claims.Username != "admin" {
				err = errors.New(" User not enough privilege")
				line = common.GetLine() - 1
				Controler.LogControl.WithFields(log.Fields{
					"Function": Function,
					"comment":  "L" + strconv.Itoa(line) + " - Error token user",
					"error":    err,
				}).Error()
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			r.Header.Set("User", claims.Username)
			next.ServeHTTP(w, r)
		}
		//******************************************************************************

	})

}
