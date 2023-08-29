package login

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"strconv"
	"uploader/SQL"
	common "uploader/common"
	Ctrl "uploader/controler"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
)

var JwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	Id       int
	jwt.RegisteredClaims
}

func Signin(w http.ResponseWriter, r *http.Request) {

	Function := "[Signin]"
	var line int

	//-----------------------------Init Controler-----------------------------------
	var Controler Ctrl.ControlerStruct
	err := Controler.ControlLogAndDB(w, "Signin")
	if err != nil {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Init Controller",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//------------------------------------------------------------------------------

	//-------------------------------Decode Body------------------------------------
	var creds Credentials
	err = json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Init Controller",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//------------------------------------------------------------------------------

	//-------------------------------Verif User-------------------------------------

	// Recup all user
	MapUsers, err := SQL.SELECTAllUser(Controler.LogControl, Controler.DB)
	if err != nil {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Init Controller",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verif si le user Existe que le password est le bon
	expectedPassword := MapUsers[creds.Username].Password
	if _, ok := MapUsers[creds.Username]; !ok || expectedPassword != creds.Password {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - User not existe or Wrong pwd",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	line = common.GetLine()
	Controler.LogControl.WithFields(log.Fields{
		"Function": Function,
		"comment":  "L" + strconv.Itoa(line) + " - User  Authent !",
		"User":     creds.Username,
	}).Info()
	//------------------------------------------------------------------------------

	//-----------------------------Token Generation---------------------------------

	// date d'expiration du token
	expirationTime := time.Now().Add(5 * time.Minute)

	// create claims => condition de validité du token
	claims := &Claims{
		Username: creds.Username,
		Id:       MapUsers[creds.Username].Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// création du token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// récup du token en string
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Generate token",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//------------------------------------------------------------------------------

	//-----------------------------Envoie du token----------------------------------
	w.Header().Set("Name", "token")
	w.Header().Set("Value", tokenString)
	w.Header().Set("Expires", expirationTime.String())
	//------------------------------------------------------------------------------
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	Function := "[Refresh]"
	var line int

	//-----------------------------Init Controler-----------------------------------
	var Controler Ctrl.ControlerStruct
	err := Controler.ControlLogAndDB(w, "Refresh")
	if err != nil {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Init Controller",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//------------------------------------------------------------------------------

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//------------------------------------------------------------------------------

	//---------------------vérfication de la validité du token----------------------

	strArr := strings.Split(bearToken, " ")
	tknStr := strArr[1]
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
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

	// check que le token n'a plus qu'1 minute de validité
	if time.Until(claims.ExpiresAt.Time) > 1*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	line = common.GetLine()
	Controler.LogControl.WithFields(log.Fields{
		"Function": Function,
		"comment":  "L" + strconv.Itoa(line) + " - Token Valid !",
	}).Info()
	//------------------------------------------------------------------------------

	//-----------------------------Token Generation---------------------------------
	// re-création du token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Generate token",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//------------------------------------------------------------------------------

	//-----------------------------Envoie du token----------------------------------
	// envoie du token
	w.Header().Set("Name", "token")
	w.Header().Set("Value", tokenString)
	w.Header().Set("Expires", expirationTime.String())
	//------------------------------------------------------------------------------
}

// func Logout(w http.ResponseWriter, r *http.Request) {
// 	// immediately clear the token cookie
// 	http.SetCookie(w, &http.Cookie{
// 		Name:    "token",
// 		Expires: time.Now(),
// 	})
// }
