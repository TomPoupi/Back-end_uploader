package SQL

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"uploader/common"

	log "github.com/sirupsen/logrus"
)

func SELECTAllVideo(logSQL *log.Logger, db *sql.DB) (map[int]common.VideoGene, error) {

	Function := "[SELECTAllVideo]"
	var line int

	var Id int
	var Name sql.NullString
	var Description sql.NullString
	var Date sql.NullString
	var Video_id sql.NullString
	var File_name sql.NullString
	var Path sql.NullString
	var Size sql.NullInt64

	MapVideo := make(map[int]common.VideoGene)

	query := "SELECT `id`,`name`,`description`,`date`,`video_id`,`file_name`,`path`,`size` FROM `project_uploader`.video_gene" +
		" INNER JOIN `video_detail` ON `object_video` = `video_detail`.`video_id`;"

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
		err = results.Scan(&Id, &Name, &Description, &Date, &Video_id, &File_name, &Path, &Size)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		OneVideoDetail := common.VideoDetail{
			Video_id:  Video_id.String,
			File_name: File_name.String,
			Path:      Path.String,
			Size:      Size.Int64,
		}
		OneVideoGene := common.VideoGene{
			Id:           Id,
			Name:         Name.String,
			Description:  Description.String,
			Date:         Date.String,
			Object_video: OneVideoDetail,
		}

		MapVideo[Id] = OneVideoGene
	}

	return MapVideo, nil
}

func SELECTOneVideo(logSQL *log.Logger, db *sql.DB, id int) (map[int]common.VideoGene, error) {

	Function := "[SELECTOneVideo]"
	var line int

	var Id int
	var Name sql.NullString
	var Description sql.NullString
	var Date sql.NullString
	var Video_id sql.NullString
	var File_name sql.NullString
	var Path sql.NullString
	var Size sql.NullInt64

	MapVideo := make(map[int]common.VideoGene)

	query := "SELECT `id`,`name`,`description`,`date`,`video_id`,`file_name`,`path`,`size` FROM `project_uploader`.video_gene" +
		" INNER JOIN `video_detail` ON `object_video` = `video_detail`.`video_id`" +
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
		err = results.Scan(&Id, &Name, &Description, &Date, &Video_id, &File_name, &Path, &Size)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		OneVideoDetail := common.VideoDetail{
			Video_id:  Video_id.String,
			File_name: File_name.String,
			Path:      Path.String,
			Size:      Size.Int64,
		}
		OneVideoGene := common.VideoGene{
			Id:           Id,
			Name:         Name.String,
			Description:  Description.String,
			Date:         Date.String,
			Object_video: OneVideoDetail,
		}

		MapVideo[Id] = OneVideoGene
	}

	return MapVideo, nil
}

func SELECTAllId(logSQL *log.Logger, db *sql.DB) (map[int]string, error) {

	Function := "[SELECTAllId]"

	var line int

	var id sql.NullInt64

	MapVideoID := make(map[int]string)

	query := "SELECT DISTINCT(id) FROM `project_uploader`.video_gene;"

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
		err = results.Scan(&id)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		MapVideoID[int(id.Int64)] = "id"
	}

	return MapVideoID, nil
}

func SELECTAllVideoID(logSQL *log.Logger, db *sql.DB) (map[int]string, error) {

	Function := "[SELECTAllVideoID]"

	var line int

	var object_video sql.NullString

	MapVideoID := make(map[int]string)

	query := "SELECT DISTINCT(object_video) FROM `project_uploader`.video_gene;"

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

func INSERTNewVideo(logSQL *log.Logger, db *sql.DB, Video common.VideoGene) error {

	Function := "[INSERTNewVideo]"

	var line int

	query1 := "INSERT INTO `project_uploader`.video_gene (`id`,`name`,`description`,`date`, `object_video`) VALUE (?,?,?,NOW(),?);"
	query2 := "INSERT INTO `project_uploader`.video_detail (`video_id`,`file_name`,`path`,`size`) VALUE (?,?,?,?);"

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

	_, err = stmt.Exec(Video.Id, Video.Name, Video.Description, Video.Object_video.Video_id)
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

	_, err = stmt2.Exec(Video.Object_video.Video_id, Video.Object_video.File_name, Video.Object_video.Path, Video.Object_video.Size)
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

func UPDATEOneVideo(logSQL *log.Logger, db *sql.DB, Body common.VideoGene, id int) error {

	Function := "[UPDATEOneVideo]"

	var line int

	qParts := make([]string, 0)
	args := make([]interface{}, 0)
	var allqpart string

	query := "UPDATE `project_uploader`.`video_gene` SET "

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

func DELETEOneUpload(logSQL *log.Logger, db *sql.DB, upload common.VideoGene) error {

	Function := "[DELETEOneUpload]"

	var line int

	query1 := "DELETE FROM `project_uploader`.`video_detail` WHERE  `video_id` = ?;"
	query2 := "DELETE FROM `project_uploader`.`video_gene` WHERE  `id` = ?;"

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

func DELETEAllUpload(logSQL *log.Logger, db *sql.DB) error {

	Function := "[DELETEAllUpload]"

	var line int

	query1 := "DELETE FROM `project_uploader`.`video_detail`;"
	query2 := "DELETE FROM `project_uploader`.`video_gene`;"

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
