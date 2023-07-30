package SQL

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func ConnProjectUploader() (*sql.DB, error) {

	dbuser := "tom"
	dbpassword := "tom"
	dburl := "localhost"
	dbname := "projet_uploader"

	fonction := "[main]"

	db, err := sql.Open("mysql", dbuser+":"+dbpassword+"@tcp("+dburl+":3306)/"+dbname+"?timeout=5s&parseTime=true")
	if err != nil {
		fmt.Println(fonction+" - line 21 : Failed connect Database", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(fonction+" - line 26 : Failed to Ping Database", err)
		return nil, err
	}
	fmt.Println("Database connection done")

	return db, nil

}
