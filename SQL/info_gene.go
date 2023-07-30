package SQL

import (
	"database/sql"
	"fmt"
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
