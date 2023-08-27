package video

import (
	"net/http"
	"os"
	"strconv"
	"uploader/SQL"
	"uploader/common"

	log "github.com/sirupsen/logrus"
)

func OperationAllData(w http.ResponseWriter, r *http.Request) {

	Function := "[OperationAllData]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	err := Controler.ControlLogAndDB(w, "OperationAllData")
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
			"comment":  "L" + strconv.Itoa(line) + " - OperationAllData GET Done",
		}).Info()

		return

		//--------------------------------------------------------------------------
	}
	//**************************************************************************

	//****************************** DELETE ************************************
	if r.Method == "DELETE" {

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

		//------------------------------DELETE File---------------------------------
		for id := range mapVideo {

			// check if exist
			if _, err := os.Stat(mapVideo[id].Object_video.Path); err == nil {
				// remove it
				err = os.Remove(mapVideo[id].Object_video.Path)
				if err != nil {
					line = common.GetLine() - 1
					Controler.LogControl.WithFields(log.Fields{
						"Function": Function,
						"comment":  "L" + strconv.Itoa(line) + " - Error on Remove file",
						"error":    err,
					}).Error()
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			line = common.GetLine()
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Delete File : " + mapVideo[id].Object_video.File_name + " Done",
			}).Info()
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Delete All File Done",
		}).Info()
		//--------------------------------------------------------------------------

		//-----------------------------DELETE All Object--------------------------------
		err = SQL.DeleteAllUpload(Controler.LogControl, Controler.DB)
		if err != nil {

			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on func SQL.DeleteOneVideo",
				"error":    err,
			}).Error()

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Delete All Done",
		}).Info()
		//--------------------------------------------------------------------------

		//-------------------------- Body Response ---------------------------------

		common.JSONresponse(Controler.LogControl, w, 200, "Delete All Done")

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - OperationAllData DELETE Done",
		}).Info()
		//--------------------------------------------------------------------------

	}
	//**************************************************************************

}
