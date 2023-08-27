package video

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"time"
	"uploader/SQL"
	"uploader/common"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func OperationOneData(w http.ResponseWriter, r *http.Request) {

	Function := "[OperationOneData]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	err := Controler.ControlLogAndDB(w, "OperationOneData")
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

	//--------------------------------END Init---------------------------------------

	//********************************* GET ************************************
	if r.Method == "GET" {

		//-------------------------------Recup varURL----------------------------------

		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Convert String to Int",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Get One Data with id : " + mux.Vars(r)["id"],
		}).Info()

		//--------------------------------------------------------------------------

		//-------------------------- GET One video ---------------------------------
		mapVideo, err := SQL.GetOneVideo(Controler.LogControl, Controler.DB, id)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on GetOneVideo",
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

	//********************************* PUT ************************************
	if r.Method == "PUT" {
		//--------------------------------Recup Body----------------------------------

		var Body common.Upload

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Recupération du body",
		}).Info()

		err = json.NewDecoder(r.Body).Decode(&Body)
		switch {

		//test si le body est vide
		case err == io.EOF:
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Resquest Body is empty",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		// s'il y a un probleme de decode
		case err != nil:
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on decode Body",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		//si le body contient un object vide
		if reflect.ValueOf(Body).IsZero() {
			err = errors.New("all fields empty")
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - format body error",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//--------------------------------------------------------------------------

		//-------------------------------Recup varURL-------------------------------

		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Convert String to Int",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Get One Data with id : " + mux.Vars(r)["id"],
		}).Info()

		//--------------------------------------------------------------------------

		//-------------------------------Update Data-------------------------------

		err = SQL.UpdateOneUpload(Controler.LogControl, Controler.DB, Body, id)
		if err != nil {

			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on func SQL.UpdateOneUpload",
				"error":    err,
			}).Error()

			if err.Error() == "Bad Format Body" {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//--------------------------------------------------------------------------

		//-------------------------- Body Response ---------------------------------

		common.JSONresponse(Controler.LogControl, w, 200, "Update Done")

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - OperationOneData PUT Done",
		}).Info()
		//--------------------------------------------------------------------------
	}
	//**************************************************************************

	//****************************** DELETE ************************************
	if r.Method == "DELETE" {

		//-------------------------------Recup varURL-------------------------------

		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Convert String to Int",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Get One Data with id : " + mux.Vars(r)["id"],
		}).Info()

		//--------------------------------------------------------------------------

		//------------------------------GET Object Video---------------------------

		mapVideo, err := SQL.GetOneVideo(Controler.LogControl, Controler.DB, id)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on GetOneVideo",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// check if exist
		if len(mapVideo) == 0 {

			err = errors.New("Video do not exist")
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Video do not exist",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		//--------------------------------------------------------------------------

		//------------------------------DELETE File---------------------------------

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
			"comment":  "L" + strconv.Itoa(line) + " - Delete File Done",
		}).Info()

		//--------------------------------------------------------------------------

		//-----------------------------DELETE Object--------------------------------
		err = SQL.DeleteOneUpload(Controler.LogControl, Controler.DB, mapVideo[id])
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
			"comment":  "L" + strconv.Itoa(line) + " - Delete Object Done",
		}).Info()
		//--------------------------------------------------------------------------

		//-------------------------- Body Response ---------------------------------

		common.JSONresponse(Controler.LogControl, w, 200, "Delete Done")

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - OperationOneData DELETE Done",
		}).Info()
		//--------------------------------------------------------------------------

	}
	//**************************************************************************

}

func GetVideoOneData(w http.ResponseWriter, r *http.Request) {

	Function := "[GetVideoOneData]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	err := Controler.ControlLogAndDB(w, "GetVideoOneData")
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

	//--------------------------------END Init---------------------------------------

	//********************************* GET ************************************

	if r.Method == "GET" {

		//-------------------------------Recup varURL----------------------------------

		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Convert String to Int",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Get One Data with id : " + mux.Vars(r)["id"],
		}).Info()

		//-------------------------------Recup varURL----------------------------------

		//-------------------------- GET One video ---------------------------------
		mapVideo, err := SQL.GetOneVideo(Controler.LogControl, Controler.DB, id)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on GetOneVideo",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//--------------------------------------------------------------------------

		//------------------ get Object Video and expose file ---------------------
		baseName := path.Base(mapVideo[id].Object_video.Path)
		file, err := os.Open(mapVideo[id].Object_video.Path)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Open Video File",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.ServeContent(w, r, baseName, time.Unix(0, 0), file)

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Expose Video file Done",
		}).Info()
		//--------------------------------------------------------------------------
	}
	//**************************************************************************

}
