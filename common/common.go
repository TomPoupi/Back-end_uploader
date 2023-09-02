package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"math/rand"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var SecretKey = []byte("N1PCdw3M2B1TfJhoaY2mL736p2vCUc47")

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
func Encrypt(logCommon *log.Logger, plaintext string) (string, error) {

	Function := "[Encrypt]"
	var line int

	aes, err := aes.NewCipher(SecretKey)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Generaate Cipher block 1",
			"error":    err,
		}).Error()
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on GCM ",
			"error":    err,
		}).Error()
		return "", err
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Genetate random ",
			"error":    err,
		}).Error()
		return "", err
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext), nil
}

// Decrypt
func Decrypt(logCommon *log.Logger, ciphertext string) (string, error) {

	Function := "[Encrypt]"
	var line int

	aes, err := aes.NewCipher(SecretKey)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Generaate Cipher block 1",
			"error":    err,
		}).Error()
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on GCM",
			"error":    err,
		}).Error()
		return "", err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		line = GetLine() - 1
		logCommon.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Decode",
			"error":    err,
		}).Error()
		return "", err
	}

	return string(plaintext), nil
}
