package logger

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

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

func Init() {

	logOne := log.New()              //faire une nouvelle instance de log
	logOne.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	logname := "main.log"
	projectbase := FindProjectPath()
	path := projectbase + "/logs"
	filename := path + "/" + logname

	if _, err := os.Stat(path); os.IsNotExist(err) {

		err := os.Mkdir(path, 0700)
		if err != nil {
			log.Fatal("Create directory failed : ", path)
		}
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//on déclare la sortie de fichier
	logOne.SetOutput(file)

	// on formatte le fichier
	logOne.SetFormatter(&log.TextFormatter{
		DisableColors:             false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		EnvironmentOverrideColors: false,
		DisableQuote:              false,
	})

	log.RegisterExitHandler(func() {
		if file == nil {
			return
		}
		file.Close()
	})

}

func LogSwitch(logname string) (*log.Logger, *os.File) {
	//On définit le nom du fichier de log

	if logname == "" {
		logname = "main.log"
	}
	logOne := log.New()              //faire une nouvelle instance de log
	logOne.SetOutput(ioutil.Discard) // Send all logs to nowhere by default

	projectbase := FindProjectPath()
	path := projectbase + "/logs"
	filename := path + "/" + logname

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//on déclare la sortie de fichier
	logOne.SetOutput(file)

	// on formatte le fichier
	logOne.SetFormatter(&log.TextFormatter{
		DisableColors:             false,
		FullTimestamp:             true,
		TimestampFormat:           "2006-01-02 15:04:05",
		EnvironmentOverrideColors: false,
		DisableQuote:              false,
	})

	log.RegisterExitHandler(func() {
		if file == nil {
			return
		}
		file.Close()
	})

	return logOne, file
}
