package SQL

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	common "uploader/common"

	log "github.com/sirupsen/logrus"
)

func SELECTAllUserByUsername(logSQL *log.Logger, db *sql.DB) (map[string]common.Users, error) {

	Function := "[SELECTAllUserByUsername]"
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

func SELECTAllIdUsers(logSQL *log.Logger, db *sql.DB) (map[int]string, error) {

	Function := "[SELECTAllIdUsers]"

	var line int

	var id sql.NullInt64

	MapUsersID := make(map[int]string)

	query := "SELECT DISTINCT(id) FROM `project_uploader`.users;"

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

		MapUsersID[int(id.Int64)] = "id"
	}

	return MapUsersID, nil
}

func INSERTNewUser(logSQL *log.Logger, db *sql.DB, User common.Users) error {

	Function := "[INSERTNewUser]"

	var line int

	query := "INSERT INTO `project_uploader`.users (`id`,`username`,`password`,`level`) VALUE (?,?,?,?);"

	stmt, err := db.Prepare(query)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Prepare INSERT",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(User.Id, User.Username, User.Password, User.Level)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on db.Exec INSERT",
			"error":    err,
		}).Error()
		return err
	}

	return nil
}

func SELECTAllUsers(logSQL *log.Logger, db *sql.DB) (map[int]common.Users, error) {

	Function := "[SELECTAllUsers]"

	var line int

	var Id sql.NullInt64
	var Username sql.NullString
	var Password sql.NullString
	var Level sql.NullInt64

	MapUsers := make(map[int]common.Users)

	query := "SELECT id, username, password, level FROM `project_uploader`.users;"

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
		err = results.Scan(&Id, &Username, &Password, &Level)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}

		OneUser := common.Users{
			Id:       int(Id.Int64),
			Username: Username.String,
			Password: Password.String,
			Level:    int(Level.Int64),
		}

		MapUsers[int(Id.Int64)] = OneUser
	}

	return MapUsers, nil
}

func SELECTOneUser(logSQL *log.Logger, db *sql.DB, id int) (map[int]common.Users, error) {

	Function := "[SELECTOneUser]"

	var line int

	var Id sql.NullInt64
	var Username sql.NullString
	var Password sql.NullString
	var Level sql.NullInt64

	MapUsers := make(map[int]common.Users)

	query := "SELECT id, username, password, level FROM `project_uploader`.users WHERE id = ?;"

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
		err = results.Scan(&Id, &Username, &Password, &Level)
		if err != nil {
			line = common.GetLine() - 1
			logSQL.WithFields(log.Fields{
				"Function": Function,
				"comment":  "L" + strconv.Itoa(line) + " - Error on Query SCAN",
				"error":    err,
			}).Error()
			return nil, err
		}
		OneUser := common.Users{
			Id:       int(Id.Int64),
			Username: Username.String,
			Password: Password.String,
			Level:    int(Level.Int64),
		}

		MapUsers[int(Id.Int64)] = OneUser
	}

	return MapUsers, nil
}

func UPDATEOneUser(logSQL *log.Logger, db *sql.DB, Body common.Users, id int) error {

	Function := "[UPDATEOneUser]"

	var line int

	qParts := make([]string, 0)
	args := make([]interface{}, 0)
	var allqpart string

	query := "UPDATE `project_uploader`.`users` SET "

	if Body.Username != "" {
		qParts = append(qParts, "`username` = ?")
		args = append(args, Body.Username)
	}
	if Body.Password != "" {
		qParts = append(qParts, "`password` = ?")
		args = append(args, Body.Password)
	}
	if Body.Level != 0 {
		qParts = append(qParts, "`level` = ?")
		args = append(args, Body.Level)
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

func DELETEOneUser(logSQL *log.Logger, db *sql.DB, User common.Users) error {

	Function := "[DELETEOneUser]"

	var line int

	query := "DELETE FROM `project_uploader`.`users` WHERE  `id` = ?;"

	stmt, err := db.Prepare(query)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Prepare ",
			"error":    err,
		}).Error()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(User.Id)
	if err != nil {
		line = common.GetLine() - 1
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Error on Exec ",
			"error":    err,
		}).Error()
		return err
	}

	return nil
}
