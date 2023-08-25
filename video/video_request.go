package video

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	SQL "uploader/SQL"
	"uploader/common"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

type myBody struct {
	name        string
	description string
}

func GetAllData(w http.ResponseWriter, r *http.Request) {

	Function := "[GetAllData]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	err := Controler.ControlLogAndDB(w, Function)
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

	if r.Method == "GET" {

		mapVideo := make(map[int]common.Upload)
		mapVideo, err := SQL.GetVideo(Controler.DB, mapVideo)
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

		response, err := json.MarshalIndent(mapVideo, "", "\t")
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on convertiting data to json",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Write Body",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - GetAllData Done",
		}).Info()

	}

}

func GetOneData(w http.ResponseWriter, r *http.Request) {

	Function := "[GetOneData]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	Controler.ControlLogAndDB(w, Function)

	//--------------------------------END Init---------------------------------------

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

		mapVideo := make(map[int]common.Upload)
		mapVideo, err = SQL.GetOneVideo(Controler.DB, mapVideo, id)
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

		response, err := json.MarshalIndent(mapVideo, "", "\t")
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on convertiting data to json",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(response)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Write Body",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - GetOneData Done",
		}).Info()

	}

}

func GetVideoOneData(w http.ResponseWriter, r *http.Request) {

	Function := "[GetVideoOneData]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	Controler.ControlLogAndDB(w, Function)

	//--------------------------------END Init---------------------------------------

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

		mapVideo := make(map[int]common.Upload)
		mapVideo, err = SQL.GetOneVideo(Controler.DB, mapVideo, id)
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

		// get Object Video and expose file
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

	}

}

func UploadVideo(w http.ResponseWriter, r *http.Request) {

	Function := "[UploadVideo]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	Controler.ControlLogAndDB(w, Function)

	//--------------------------------END Init---------------------------------------

	if r.Method == "POST" {

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Upload Begin",
		}).Info()

		//-----------------------------Récup Body-----------------------------------

		//l'ordre est important et il est important de cloner les response pour les utiliser plusieurs fois
		r1 := r
		r2 := r
		r3 := r
		// Parse our multipart form, 200 << 20 specifies a maximum
		// upload of 200 MB files.
		r1.ParseMultipartForm(200 << 20)
		r2.ParseForm()
		r3.ParseForm()
		// FormFile returns the first file for the given key `myFile`
		// it also returns the FileHeader so we can get the Filename,
		// the Header and the size of the file

		name := r2.Form.Get("name")
		description := r3.Form.Get("description")

		file, handler, err := r1.FormFile("myFile")
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Retreiving File",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Get Video Done ",
		}).Info()

		//--------------------------------------------------------------------------

		//-------------------------Create Video File--------------------------------

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function":    Function,
			"comment":     "L" + strconv.Itoa(line) + " - Info Video : ",
			"file name":   handler.Filename,
			"file size":   handler.Size,
			"mime Header": handler.Header,
		}).Info()

		// Create a new file in the uploads directory
		f, err := os.OpenFile("./file/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error to Create File",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer f.Close()

		// Copy the contents of the file to the new file
		_, err = io.Copy(f, file)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Insert Content File",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//--------------------------------------------------------------------------

		//-----------------------------Create IdVideo-------------------------------
		MapVideoID, err := SQL.GetAllVideoID(Controler.DB)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on SQL.GetAllVideoID",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var Video_id string
		for i := 1; i < len(MapVideoID)+2; i++ {

			if _, ok := MapVideoID[i]; !ok {
				Video_id = "video_" + strconv.Itoa(i)
				break
			}
		}

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - New Id Video :" + Video_id,
		}).Info()

		//--------------------------------------------------------------------------

		//-------------------------Insert Video On DataBase-------------------------
		OneVideo := common.Video{
			Video_id:  Video_id,
			File_name: handler.Filename,
			Path:      common.FindProjectPath() + "\\file\\" + handler.Filename, // à changer pour s'adapter à env linux
		}

		OneUpload := common.Upload{
			Name:         name,
			Description:  description,
			Object_video: OneVideo,
		}

		err = SQL.PostVideo(Controler.DB, OneUpload)
		if err != nil {
			line = common.GetLine() - 1
			Controler.LogControl.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on func SQL.PostVideo",
				"error":    err,
			}).Error()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//--------------------------------------------------------------------------

	}
}
