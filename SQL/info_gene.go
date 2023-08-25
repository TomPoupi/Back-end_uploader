package SQL

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"uploader/common"
)

func GetVideo(db *sql.DB, MapVideo map[int]common.Upload) (map[int]common.Upload, error) {
	function := "[GetVideo]"
	var Id int
	var Name sql.NullString
	var Description sql.NullString
	var Date sql.NullString
	var Video_id sql.NullString
	var File_name sql.NullString
	var Path sql.NullString

	query := "SELECT `id`,`name`,`description`,`date`,`video_id`,`file_name`,`path` FROM `projet_uploader`.info_gene" +
		" INNER JOIN `video` ON `object_video` = `video`.`video_id`;"

	results, err := db.Query(query)
	if err != nil {
		fmt.Println(function, "- line 21 : error on query SELECT : ", err)
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&Id, &Name, &Description, &Date, &Video_id, &File_name, &Path)
		if err != nil {
			fmt.Println(function, "- line 29 : error on scan SELECT : ", err)
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

func GetOneVideo(db *sql.DB, MapVideo map[int]common.Upload, id int) (map[int]common.Upload, error) {
	function := "[GetData]"
	var Id int
	var Name sql.NullString
	var Description sql.NullString
	var Date sql.NullString
	var Video_id sql.NullString
	var File_name sql.NullString
	var Path sql.NullString

	query := "SELECT `id`,`name`,`description`,`date`,`video_id`,`file_name`,`path` FROM `projet_uploader`.info_gene" +
		" INNER JOIN `video` ON `object_video` = `video`.`video_id`" +
		" WHERE `id` = ?;"

	results, err := db.Query(query, id)
	if err != nil {
		fmt.Println(function, "- line 21 : error on query SELECT : ", err)
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&Id, &Name, &Description, &Date, &Video_id, &File_name, &Path)
		if err != nil {
			fmt.Println(function, "- line 29 : error on scan SELECT : ", err)
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

func GetAllVideoID(db *sql.DB) (map[int]string, error) {

	function := "[GetAllVideoID]"

	var object_video sql.NullString

	MapVideoID := make(map[int]string)

	query := "SELECT DISTINCT(object_video) FROM `projet_uploader`.info_gene;"

	results, err := db.Query(query)
	if err != nil {
		fmt.Println(function, "- line 21 : error on query SELECT : ", err)
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&object_video)
		if err != nil {
			fmt.Println(function, "- line 29 : error on scan SELECT : ", err)
			return nil, err
		}

		if object_video.String == "" {
			err = errors.New("video_id ne peut être vide")
			fmt.Println("Get video_id error : ", err)
			return nil, err
		}

		id := strings.Split(object_video.String, "_")
		idInt, _ := strconv.Atoi(id[1])
		MapVideoID[idInt] = object_video.String
	}

	return MapVideoID, nil
}

func PostVideo(db *sql.DB, Video common.Upload) error {

	function := "[GetData]"

	query1 := "INSERT INTO `projet_uploader`.info_gene (`name`,`description`,`date`, `object_video`) VALUE (?,?,NOW(),?);"
	query2 := "INSERT INTO `projet_uploader`.video (`video_id`,`file_name`,`path`) VALUE (?,?,?);"

	stmt, err := db.Prepare(query1)
	if err != nil {
		fmt.Println(function+"- line 149 : error on db.Prepare INSERT 1  : ", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(Video.Name, Video.Description, Video.Object_video.Video_id)
	if err != nil {
		fmt.Println(function+"- line 154 : error on db.Exec INSERT 1 : ", err)
		return err
	}

	stmt2, err := db.Prepare(query2)
	if err != nil {
		fmt.Println(function+"- line 160 : error on db.Prepare INSERT 2 : ", err)
		return err
	}
	defer stmt2.Close()

	_, err = stmt2.Exec(Video.Object_video.Video_id, Video.Object_video.File_name, Video.Object_video.Path)
	if err != nil {
		fmt.Println(function+"- line 167 : error on db.Exec INSERT 2 : ", err)
		return err
	}

	return nil
}
