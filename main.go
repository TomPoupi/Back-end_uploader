package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"uploader/SQL"
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

	r.Handle("/login", authMiddlerware(http.HandlerFunc(login.Signin), OpenAccess)).Methods("POST")
	r.Handle("/refreshToken", authMiddlerware(http.HandlerFunc(login.Refresh), OpenAccess)).Methods("POST")

	r.Handle("/create_user", authMiddlerware(http.HandlerFunc(login.CreateLogin), RestrictedAccess)).Methods("POST")
	r.Handle("/user", authMiddlerware(http.HandlerFunc(login.OperationAllUsers), RestrictedAccess)).Methods("GET")
	r.Handle("/user/{id}", authMiddlerware(http.HandlerFunc(login.OperationOneUser), RestrictedAccess)).Methods("GET")
	r.Handle("/user/{id}", authMiddlerware(http.HandlerFunc(login.OperationOneUser), RestrictedAccess)).Methods("PUT")
	r.Handle("/user/{id}", authMiddlerware(http.HandlerFunc(login.OperationOneUser), RestrictedAccess)).Methods("DELETE")

	r.Handle("/video", authMiddlerware(http.HandlerFunc(video.OperationAllVideo), OpenAccess)).Methods("GET")
	r.Handle("/video", authMiddlerware(http.HandlerFunc(video.OperationAllVideo), RestrictedAccess)).Methods("DELETE")
	r.Handle("/video/{id}", authMiddlerware(http.HandlerFunc(video.OperationOneVideo), OpenAccess)).Methods("GET")
	r.Handle("/video/{id}", authMiddlerware(http.HandlerFunc(video.OperationOneVideo), RestrictedAccess)).Methods("PUT")
	r.Handle("/video/{id}", authMiddlerware(http.HandlerFunc(video.OperationOneVideo), RestrictedAccess)).Methods("DELETE")
	r.Handle("/video/{id}/file", authMiddlerware(http.HandlerFunc(video.GetOneVideoFile), OpenAccess)).Methods("GET")
	r.Handle("/upload_video", authMiddlerware(http.HandlerFunc(video.UploadVideo), RestrictedAccess)).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()

}

func Home(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusTeapot)
	w.Write([]byte("Hello , Welcome in Backend uploader "))

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
			claims, err := login.VerifyValideTkn(Controler.LogControl, bearToken)
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
				if err.Error() == "Invalide Token" {
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

			line = common.GetLine()
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Token Valid !",
			}).Info()
			//------------------------------------------------------------------------------

			//---------------vérfication de niveau d'accès de l'utilisateur-----------------

			MapUsers, err := SQL.SELECTOneUser(Controler.LogControl, Controler.DB, claims.Id)
			if err != nil {
				Controler.LogControl.WithFields(log.Fields{
					"Function": Function,
					"comment":  "L" + strconv.Itoa(line) + " - Error on SQL.SELECTOneUser",
					"error":    err,
				}).Error()
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if MapUsers[claims.Id].Level < 100 {
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

			// if claims.Username != "admin" {
			// 	err = errors.New(" User not enough privilege")
			// 	line = common.GetLine() - 1
			// 	Controler.LogControl.WithFields(log.Fields{
			// 		"Function": Function,
			// 		"comment":  "L" + strconv.Itoa(line) + " - Error token user",
			// 		"error":    err,
			// 	}).Error()
			// 	http.Error(w, err.Error(), http.StatusUnauthorized)
			// 	return
			// }

			//------------------------------------------------------------------------------

			r.Header.Set("User", MapUsers[claims.Id].Username)
			next.ServeHTTP(w, r)
		}
		//******************************************************************************

	})

}
