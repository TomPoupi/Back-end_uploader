package SQL

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"uploader/common"

	log "github.com/sirupsen/logrus"
)

func GetVideo(logSQL *log.Logger, db *sql.DB) (map[int]common.Upload, error) {

	Function := "[GetVideo]"
	var line int

	var Id int
	var Name sql.NullString
	var Description sql.NullString
	var Date sql.NullString
	var Video_id sql.NullString
	var File_name sql.NullString
	var Path sql.NullString

	MapVideo := make(map[int]common.Upload)

	query := "SELECT `id`,`name`,`description`,`date`,`video_id`,`file_name`,`path` FROM `projet_uploader`.info_gene" +
		" INNER JOIN `video` ON `object_video` = `video`.`video_id`;"

	results, err := db.Query(query)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Query SELECT",
			"error":    err,
		}).Error()
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&Id, &Name, &Description, &Date, &Video_id, &File_name, &Path)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		OneVideo := common.Video{
			Video_id:  Video_id.String,
			File_name: File_name.String,
			Path:      Path.String,
		}
		OneUpload := common.Upload{
			Id:           Id,
			Name:         Name.String,
			Description:  Description.String,
			Date:         Date.String,
			Object_video: OneVideo,
		}

		MapVideo[Id] = OneUpload
	}

	return MapVideo, nil
}

func GetOneVideo(logSQL *log.Logger, db *sql.DB, id int) (map[int]common.Upload, error) {

	Function := "[GetOneVideo]"
	var line int

	var Id int
	var Name sql.NullString
	var Description sql.NullString
	var Date sql.NullString
	var Video_id sql.NullString
	var File_name sql.NullString
	var Path sql.NullString

	MapVideo := make(map[int]common.Upload)

	query := "SELECT `id`,`name`,`description`,`date`,`video_id`,`file_name`,`path` FROM `projet_uploader`.info_gene" +
		" INNER JOIN `video` ON `object_video` = `video`.`video_id`" +
		" WHERE `id` = ?;"

	results, err := db.Query(query, id)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Query SELECT",
			"error":    err,
		}).Error()
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&Id, &Name, &Description, &Date, &Video_id, &File_name, &Path)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		OneVideo := common.Video{
			Video_id:  Video_id.String,
			File_name: File_name.String,
			Path:      Path.String,
		}
		OneUpload := common.Upload{
			Id:           Id,
			Name:         Name.String,
			Description:  Description.String,
			Date:         Date.String,
			Object_video: OneVideo,
		}

		MapVideo[Id] = OneUpload
	}

	return MapVideo, nil
}

func GetAllVideoID(logSQL *log.Logger, db *sql.DB) (map[int]string, error) {

	Function := "[GetAllVideoID]"

	var line int

	var object_video sql.NullString

	MapVideoID := make(map[int]string)

	query := "SELECT DISTINCT(object_video) FROM `projet_uploader`.info_gene;"

	results, err := db.Query(query)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Query SELECT",
			"error":    err,
		}).Error()
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&object_video)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		if object_video.String == "" {
			line = common.GetLine() - 1
			err = errors.New("video_id ne peut être vide")
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - get video_id error",
				"error":    err,
			}).Error()
			return nil, err
		}

		id := strings.Split(object_video.String, "_")
		idInt, _ := strconv.Atoi(id[1])
		MapVideoID[idInt] = object_video.String
	}

	return MapVideoID, nil
}

func PostVideo(logSQL *log.Logger, db *sql.DB, Video common.Upload) error {

	Function := "[GetData]"

	var line int

	query1 := "INSERT INTO `projet_uploader`.info_gene (`name`,`description`,`date`, `object_video`) VALUE (?,?,NOW(),?);"
	query2 := "INSERT INTO `projet_uploader`.video (`video_id`,`file_name`,`path`) VALUE (?,?,?);"

	stmt, err := db.Prepare(query1)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Prepare INSERT 1",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(Video.Name, Video.Description, Video.Object_video.Video_id)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Exec INSERT 1",
			"error":    err,
		}).Error()
		return err
	}

	stmt2, err := db.Prepare(query2)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Prepare INSERT 2",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt2.Close()

	_, err = stmt2.Exec(Video.Object_video.Video_id, Video.Object_video.File_name, Video.Object_video.Path)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Exec INSERT 2",
			"error":    err,
		}).Error()
		return err
	}

	return nil
}

func UpdateOneUpload(logSQL *log.Logger, db *sql.DB, Body common.Upload, id int) error {

	Function := "[UpdateOneUpload]"

	var line int

	qParts := make([]string, 0)
	args := make([]interface{}, 0)
	var allqpart string

	query := "UPDATE `projet_uploader`.`info_gene` SET "

	if Body.Name != "" {
		qParts = append(qParts, "`name` = ?")
		args = append(args, Body.Name)
	}
	if Body.Description != "" {
		qParts = append(qParts, "`description` = ?")
		args = append(args, Body.Description)
	}

	fin := " WHERE `id` = ?;"
	args = append(args, id)

	if len(qParts) == 0 || len(args) == 0 {
		err := errors.New("Bad Format Body")
		return err
	}

	allqpart += strings.Join(qParts, ",")
	query = query + allqpart + fin

	stmt, err := db.Prepare(query)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Prepare ",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Exec ",
			"error":    err,
		}).Error()
		return err
	}

	return nil
}

func DeleteOneUpload(logSQL *log.Logger, db *sql.DB, upload common.Upload) error {

	Function := "[DeleteOneUpload]"

	var line int

	query1 := "DELETE FROM `projet_uploader`.`video` WHERE  `video_id` = ?;"
	query2 := "DELETE FROM `projet_uploader`.`info_gene` WHERE  `id` = ?;"

	stmt1, err := db.Prepare(query1)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Prepare 1",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt1.Close()

	_, err = stmt1.Exec(upload.Object_video.Video_id)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Exec 1",
			"error":    err,
		}).Error()
		return err
	}

	stmt2, err := db.Prepare(query2)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Prepare 2",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt2.Close()

	_, err = stmt2.Exec(upload.Id)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Exec 2",
			"error":    err,
		}).Error()
		return err
	}

	return nil
}

func DeleteAllUpload(logSQL *log.Logger, db *sql.DB) error {

	Function := "[DeleteAllUpload]"

	var line int

	query1 := "DELETE FROM `projet_uploader`.`video`;"
	query2 := "DELETE FROM `projet_uploader`.`info_gene`;"

	stmt1, err := db.Prepare(query1)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Prepare 1",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt1.Close()

	_, err = stmt1.Exec()
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Exec 1",
			"error":    err,
		}).Error()
		return err
	}

	stmt2, err := db.Prepare(query2)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Prepare 2",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt2.Close()

	_, err = stmt2.Exec()
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Exec 2",
			"error":    err,
		}).Error()
		return err
	}

	return nil
}
