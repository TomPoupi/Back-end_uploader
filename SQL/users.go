package SQL

import (
	"database/sql"
	"strconv"
	common "uploader/common"

	log "github.com/sirupsen/logrus"
)

func SELECTAllUser(logSQL *log.Logger, db *sql.DB) (map[string]common.Users, error) {

	Function := "[SELECTAllUser]"
	var line int

	var Id int
	var Username sql.NullString
	var Password sql.NullString

	MapUsers := make(map[string]common.Users)

	query := "SELECT `id`,`username`,`password` FROM `project_uploader`.users ;"

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
		err = results.Scan(&Id, &Username, &Password)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		OneUsers := common.Users{
			Id:       Id,
			Username: Username.String,
			Password: Password.String,
		}

		MapUsers[Username.String] = OneUsers
	}

	return MapUsers, nil
}
