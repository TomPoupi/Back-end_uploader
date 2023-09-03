package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// SecretKey len = 24
var SecretKey = []byte("abc&1*~#^2^#s1^=)^^7%b34")

// suite de Byte utilisé pour l'encryptage
var Bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func FindProjectPath() string {

	var (
		_, b, _, _  = runtime.Caller(0)
		projectbase = filepath.Join(filepath.Dir(b), "..")
	)

	//m1 := regexp.MustCompile(`\\`)
	//projectbaseclean := m1.ReplaceAllString(projectbase, "/")

	return projectbase
}

func GetLine() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}

func CreatePath(folders []string) string {

	var path string
	if runtime.GOOS == "windows" {
		path = "\\" + strings.Join(folders, "\\") + "\\"
	} else {
		path = "/" + strings.Join(folders, "/") + "/"
	}
	return path
}

func JSONresponse(logCommon *log.Logger, w http.ResponseWriter, StatusCode int, payload interface{}) {

	Function := "[JSONresponse]"
	var line int

	response, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on convertiting data to json",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(StatusCode)
	_, err = w.Write(response)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Write Body",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Encrypt

func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func Decode(logCommon *log.Logger, s string) ([]byte, error) {

	Function := "[Decode]"
	var line int

	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on DecodeString",
			"error":    err,
		}).Error()
		return nil, err
	}
	return data, nil
}

func Encrypt(logCommon *log.Logger, text string, MySecret string) (string, error) {

	Function := "[Encrypt]"
	var line int

	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Encode Key",
			"error":    err,
		}).Error()
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, Bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return Encode(cipherText), nil
}

// Decrypt
// Decrypt method is to extract back the encrypted text
func Decrypt(logCommon *log.Logger, text string, MySecret string) (string, error) {

	Function := "[Decrypt]"
	var line int

	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Encode Key",
			"error":    err,
		}).Error()
		return "", err
	}

	cipherText, err := Decode(logCommon, text)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Decode",
			"error":    err,
		}).Error()
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, Bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

func RecupSecretKey(logCommon *log.Logger) (SecretKeyStruct, error) {

	Function := "[RecupSecretKey]"
	var line int

	var Secret SecretKeyStruct

	content, err := ioutil.ReadFile("./common/secret.json")
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Open file",
			"error":    err,
		}).Error()
		return Secret, err
	}

	err = json.Unmarshal(content, &Secret)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Unmarshal",
			"error":    err,
		}).Error()

		return Secret, err

	}

	return Secret, nil
}
