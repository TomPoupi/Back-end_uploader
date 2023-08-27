package video

import (
	"net/http"
	"strconv"
	"uploader/SQL"
	"uploader/common"

	log "github.com/sirupsen/logrus"
)

func GetAllData(w http.ResponseWriter, r *http.Request) {

	Function := "[GetAllData]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	err := Controler.ControlLogAndDB(w, "GetAllData")
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

		mapVideo, err := SQL.GetVideo(Controler.LogControl, Controler.DB)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error Get Video",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//--------------------------------------------------------------------------

		//-------------------------- Body Response ---------------------------------
		common.JSONresponse(Controler.LogControl, w, 200, mapVideo)

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - OperationOneData GET Done",
		}).Info()

		return

		//--------------------------------------------------------------------------
	}
	//**************************************************************************

}
