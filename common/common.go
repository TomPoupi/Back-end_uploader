package common

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

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
