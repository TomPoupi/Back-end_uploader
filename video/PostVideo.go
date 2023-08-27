package video

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"uploader/SQL"
	"uploader/common"

	log "github.com/sirupsen/logrus"
)

type myBody struct {
	name        string
	description string
}

func UploadVideo(w http.ResponseWriter, r *http.Request) {

	Function := "[UploadVideo]"
	var line int

	//-----------------------------Init Controler-----------------------------------

	var Controler ControlerStruct
	err := Controler.ControlLogAndDB(w, "UploadVideo")
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
		MapVideoID, err := SQL.GetAllVideoID(Controler.LogControl, Controler.DB)
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
		folders := make([]string, 0)
		folders = append(folders, "file")
		OneVideo := common.Video{
			Video_id:  Video_id,
			File_name: handler.Filename,
			Path:      common.FindProjectPath() + common.CreatePath(folders) + handler.Filename, // à changer pour s'adapter à env linux
		}

		OneUpload := common.Upload{
			Name:         name,
			Description:  description,
			Object_video: OneVideo,
		}

		err = SQL.PostVideo(Controler.LogControl, Controler.DB, OneUpload)
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

		//-------------------------- Body Response ---------------------------------

		common.JSONresponse(Controler.LogControl, w, 200, "Insert Done")

		line = common.GetLine()
		Controler.LogControl.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - UploadVideo Done",
		}).Info()

		//------------------------------------------------------------------------

	}
}
