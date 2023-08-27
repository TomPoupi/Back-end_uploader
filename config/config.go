package config

import (
	"errors"
	"os"
	"strconv"
	"uploader/common"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Init(logInit *log.Logger) {

	Function := "[Init]"
	var line int

	line = common.GetLine()
	logInit.WithFields(log.Fields{
		"Function": Function,
		"comment":  "L" + strconv.Itoa(line) + " - Reading OSENV",
	}).Info()

	OSenv := os.Getenv("OSENV")

	line = common.GetLine()
	logInit.WithFields(log.Fields{
		"Function": Function,
		"comment":  "L" + strconv.Itoa(line) + " - OSENV = " + OSenv,
	}).Info()

	if OSenv == "devlocal" {
		//find filename
		viper.SetConfigName("localhostconfig")
		line = common.GetLine()
		logInit.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - on utilise le fichier de config \"localhostconfig\" ",
		}).Info()

	} else {

		err := errors.New("No OSENV found")
		line = common.GetLine()
		logInit.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - problème de lecture de la var d'environement ",
			"error":    err,
		}).Error()
		return
	}

	//find folder
	projectbase := common.FindProjectPath()
	viper.AddConfigPath(projectbase + "/config/")

	//find extention
	viper.SetConfigType("yml")

	//reading object
	if err := viper.ReadInConfig(); err != nil {
		line = common.GetLine()
		logInit.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Can't read object ",
			"error":    err,
		}).Error()
		return
	}
}
