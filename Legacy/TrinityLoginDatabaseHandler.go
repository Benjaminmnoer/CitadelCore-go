package Legacy

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const setQuery = "UPDATE account SET session_key_auth = ?, last_login = NOW() WHERE username = ?"
const getQuery = "SELECT session_key_auth FROM account WHERE username = ?"

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

	_, queryErr := db.Query(setQuery, authkey, username)
	if queryErr != nil {
		return fmt.Errorf("Query error, %s\n", queryErr)
	}

	db.Close()
	return nil
}

func GetSessionKey(username string) ([]byte, error) {
	config := mysql.NewConfig()
	config.User = "citadelcore"
	config.Passwd = "citadelcore"
	config.Addr = "127.0.0.1:3306"
	config.DBName = "auth"

	var err error
	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("DB error, %s\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Db error, %s\n", err)
	}

	rows, err := db.Query(getQuery, username)
	if err != nil {
		return nil, fmt.Errorf("Query error, %s\n", err)
	}
	defer rows.Close()

	rows.Next()
	result := make([]byte, 20)
	err = rows.Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("Scan error, %s\n", err)
	}

	return result, nil
}
