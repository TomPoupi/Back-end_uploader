package common

import "github.com/golang-jwt/jwt/v4"

type VideoGene struct {
	Id           int
	Name         string
	Description  string
	Date         string
	Object_video VideoDetail
}

type VideoDetail struct {
	Video_id  string
	File_name string
	Path      string
	Size      int64
}

type Users struct {
	Id       int
	Username string `json:"username"`
	Password string `json:"password"`
	Level    int    `json:"level"`
	Key      string
}

type Claims struct {
	Username string `json:"username"`
	Id       int
	jwt.RegisteredClaims
}

type ClaimsGene struct {
	TokenUser  string
	KeyCrypted []byte
	jwt.RegisteredClaims
}

// type News struct {
// 	Headline string
// 	Body     string
// }
