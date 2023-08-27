package SQL

import (
	"database/sql"
	"strconv"
	common "uploader/common"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ConnProjectUploader(logSQL *log.Logger) (*sql.DB, error) {

	dbuser := viper.GetString("project_uploader.dbuser")
	dbpassword := viper.GetString("project_uploader.dbpassword")
	dburl := viper.GetString("project_uploader.dburl")
	dbname := viper.GetString("project_uploader.dbname")

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

	line = common.GetLine()
	logSQL.WithFields(log.Fields{
		"Function": Function,
		"comment":  "L" + strconv.Itoa(line) + " - Database connection done",
	}).Info()

	return db, nil

}
