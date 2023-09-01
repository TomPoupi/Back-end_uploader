package login

import (
	"net/http"
	"strconv"
	"uploader/SQL"
	common "uploader/common"
	Ctrl "uploader/controler"

	log "github.com/sirupsen/logrus"
)

func OperationAllUsers(w http.ResponseWriter, r *http.Request) {

	Function := "[OperationAllUsers]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler Ctrl.ControlerStruct
	err := Controler.ControlLogAndDB(w, "OperationAllUsers")
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

	//--------------------------------------------------------------------------

	//********************************* GET ************************************
	if r.Method == "GET" {

		//-------------------------- GET All video ---------------------------------

		mapUsers, err := SQL.SELECTAllUsers(Controler.LogControl, Controler.DB)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error SQL.SELECTAllUsers",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//--------------------------------------------------------------------------

		//-------------------------- Body Response ---------------------------------
		common.JSONresponse(Controler.LogControl, w, 200, mapUsers)

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - OperationAllUsers GET Done",
		}).Info()

		return

		//--------------------------------------------------------------------------
	}
	//**************************************************************************

}
