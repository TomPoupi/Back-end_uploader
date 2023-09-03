package login

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"uploader/SQL"
	common "uploader/common"
	Ctrl "uploader/controler"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func OperationOneUser(w http.ResponseWriter, r *http.Request) {

	Function := "[OperationOneUser]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler Ctrl.ControlerStruct
	err := Controler.ControlLogAndDB(w, "OperationOneUser")
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

	//-------------------------------Verify User-------------------------------------

	UserLevel, err := strconv.Atoi(r.Header.Get("UserLevel"))
	if err != nil {
		line = common.GetLine() - 1
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error Convert string to int",
			"error":    err,
		}).Error()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if UserLevel < 100 {
		if r.Header.Get("UserId") != mux.Vars(r)["id"] {
			err := errors.New("User UnAuthorized")
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - User UnAuthorized",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return

		}
	}

	//-------------------------------------------------------------------------------

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
			"comment":  "L" + strconv.Itoa(line) + " - Get One User with id : " + mux.Vars(r)["id"],
		}).Info()

		//--------------------------------------------------------------------------

		//-------------------------- GET One user ---------------------------------
		mapUsers, err := SQL.SELECTOneUser(Controler.LogControl, Controler.DB, id)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on SQL.SELECTOneUser",
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
			"comment":  "L" + strconv.Itoa(line) + " - OperationOneUser GET Done",
		}).Info()

		return

		//--------------------------------------------------------------------------
	}
	//**************************************************************************

	//********************************* PUT ************************************
	if r.Method == "PUT" {
		//--------------------------------Recup Body----------------------------------

		var Body common.Users

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
			"comment":  "L" + strconv.Itoa(line) + " - Get One User with id : " + mux.Vars(r)["id"],
		}).Info()

		//--------------------------------------------------------------------------

		//-------------------------------Update Data-------------------------------

		err = SQL.UPDATEOneUser(Controler.LogControl, Controler.DB, Body, id)
		if err != nil {

			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on func SQL.UPDATEOneUser",
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

		common.JSONresponse(Controler.LogControl, w, 200, "Update User Done")

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - OperationOneUser PUT Done",
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
			"comment":  "L" + strconv.Itoa(line) + " - Get One User with id : " + mux.Vars(r)["id"],
		}).Info()

		//--------------------------------------------------------------------------

		//------------------------------GET Object User-----------------------------

		mapUsers, err := SQL.SELECTOneUser(Controler.LogControl, Controler.DB, id)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on SQL.SELECTOneUser",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// check if exist
		if len(mapUsers) == 0 {

			err = errors.New("User do not exist")
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - User do not exist",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		// check if not admin
		if mapUsers[id].Username == "admin" {

			err = errors.New("Cannot delete User : admin")
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Wrong delete User",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return

		}

		//--------------------------------------------------------------------------

		//-----------------------------DELETE Object User---------------------------
		err = SQL.DELETEOneUser(Controler.LogControl, Controler.DB, mapUsers[id])
		if err != nil {

			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on func SQL.DELETEOneUpload",
				"error":    err,
			}).Error()

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Delete Object User Done",
		}).Info()
		//--------------------------------------------------------------------------

		//-------------------------- Body Response ---------------------------------

		common.JSONresponse(Controler.LogControl, w, 200, "Delete User Done")

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - OperationOneUser DELETE Done",
		}).Info()
		//--------------------------------------------------------------------------

	}
	//**************************************************************************

}
