package video

import (
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"uploader/SQL"
	"uploader/common"
	"uploader/logger"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type ControlerStruct struct {
	Logname    string
	LogControl *log.Logger
	DB         *sql.DB
	Logfile    *os.File
}

func (Controller *ControlerStruct) ControlLogAndDB(w http.ResponseWriter, fonction string) error {

	Operation := "ControlLogAndDB"
	var line int

	//-------------------------- Begin Log-----------------------------

	logname := fonction + ".log"
	if fonction == "" {
		logname = "main.log"
	}

	logcontrol, logfile := logger.LogSwitch(logname)

	logcontrol.WithFields(log.Fields{
		"Operation": "------------------------------------------------------------------",
	}).Info()

	logcontrol.WithFields(log.Fields{
		"Operation": "----------------Lancement de " + fonction + "---------------------",
	}).Info()

	logcontrol.WithFields(log.Fields{
		"Operation": "------------------------------------------------------------------",
	}).Info()

	//---------------------------- End Log ----------------------------

	//---------------------------- DB init ----------------------------

	db, err := SQL.ConnProjectUploader(logcontrol)
	if err != nil {
		line = common.GetLine()
		logcontrol.WithFields(log.Fields{
			"Operation": Operation,
			"comment":   "L" + strconv.Itoa(line) + " - Failed to Connect Database",
			"error":     err,
		}).Error()

		message := "Failed to Connect Database"

		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(message))
		if err != nil {
			logcontrol.WithFields(log.Fields{
				"Operation": Operation,
				"comment":   "L" + strconv.Itoa(line) + " - Failed to Write Body",
				"error":     err,
			}).Error()
			return err
		}
		return err
	}

	//---------------------------- End DB -----------------------------

	Controller.DB = db
	Controller.Logname = logname
	Controller.LogControl = logcontrol
	Controller.Logfile = logfile

	return nil

}
