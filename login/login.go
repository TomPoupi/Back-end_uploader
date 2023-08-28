package login

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var JwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader((http.StatusBadRequest))
		return
	}

	// verification user
	expectedPassword, ok := users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// date d'expiration du token
	expirationTime := time.Now().Add(5 * time.Minute)

	// create claims => condition de validité du token
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// création du token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// récup du token en string
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// envoie du token
	w.Header().Set("Name", "token")
	w.Header().Set("Value", tokenString)
	w.Header().Set("Expires", expirationTime.String())
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	// prend le token actuelle
	bearToken := r.Header.Get("Authorization")
	if bearToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// vérfication de la validité du token
	strArr := strings.Split(bearToken, " ")
	tknStr := strArr[1]
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
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

	// check que le token n'a plus qu'1 minute de validité
	if time.Until(claims.ExpiresAt.Time) > 1*time.Minute {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// re-création du token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// envoie du token
	w.Header().Set("Name", "token")
	w.Header().Set("Value", tokenString)
	w.Header().Set("Expires", expirationTime.String())
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}
