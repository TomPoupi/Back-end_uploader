package SQL

import (
	"database/sql"
	"strconv"
	common "uploader/common"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func ConnProjectUploader(logSQL *log.Logger) (*sql.DB, error) {

	dbuser := "tom"
	dbpassword := "tom"
	dburl := "localhost"
	dbname := "projet_uploader"

	Function := "[ConnProjectUploader]"
	var line int

	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dburl+":3306)/"+dbname+"?timeout=5s&parseTime=true")
	if err != nil {
		line = common.GetLine()
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Failed to Connect Database",
			"error":    err,
		}).Error()
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		line = common.GetLine()
		logSQL.WithFields(log.Fields{
			"Function": Function,
			"comment":  "L" + strconv.Itoa(line) + " - Failed to Ping Database",
			"error":    err,
		}).Error()
		return nil, err
	}

	logSQL.WithFields(log.Fields{
		"Function": Function,
		"comment":  "L" + strconv.Itoa(line) + " - Database connection done",
	}).Info()

	return db, nil

}
