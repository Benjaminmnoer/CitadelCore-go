package Legacy

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const query = "UPDATE account SET session_key_auth = ?, last_login = NOW() WHERE username = ?"

func SetSessionKey(authkey []byte, username string) error {
	config := mysql.NewConfig()
	config.User = "citadelcore"
	config.Passwd = "citadelcore"
	config.Addr = "127.0.0.1:3306"
	config.DBName = "auth"

	var err error
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return fmt.Errorf("DB error, %s\n", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("Db error, %s\n", err)
	}

	_, queryErr := db.Query(query, authkey, username)
	if queryErr != nil {
		return fmt.Errorf("Query error, %s\n", queryErr)
	}

	db.Close()
	return nil
}
