package Dbconnection

import (
	"CitadelCore/Shared/Reflection"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type MysqlDatabaseConnection struct {
	config *mysql.Config
}

func Test(username string) ([]byte, error) {
	config := mysql.NewConfig()
	config.User = "trinity"
	config.Passwd = "trinity"
	config.Addr = "127.0.0.1:3306"
	config.DBName = "auth"

	var err error
	db, err = sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Printf("DB error, %s\n", err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("DB error, %s\n", err)
		return nil, errors.New("Db error")
	}
	fmt.Println("Connected to db")

	rows, queryErr := db.Query("SELECT salt FROM account WHERE username = ?", username)
	if queryErr != nil {
		fmt.Printf("Query error, %s\n", queryErr)
		return nil, errors.New("Query error")
	}

	res := make([]byte, 32)
	rows.Scan(&res)
	db.Close()
	return res, nil
}

func (connection MysqlDatabaseConnection) ExecuteQuerySingleResult(query string, result interface{}, args ...interface{}) (interface{}, error) {
	db, err := sql.Open("mysql", connection.config.FormatDSN())
	if err != nil {
		fmt.Printf("mysql-db: error opening database connection. Error: %s\n", err)
		return nil, err
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Printf("mysql-db: error executing query. Error: %s\n", err)
		return nil, err
	}

	hasresult := rows.Next()
	if !hasresult {
		return nil, errors.New("Error calling next")
	}

	fields := Reflection.GetFields(result)
	err = rows.Scan(fields...)
	if err != nil {
		fmt.Printf("mysql-db: error reading result. Error: %s\n", err)
		return nil, err
	}

	return Reflection.CreateResultlFromFields(fields, result), nil // createResultlFromFields(fields, result), nil
}

func CreateDatabaseConnection(username string, password string, address string, database string) (MysqlDatabaseConnection, error) {
	config := mysql.NewConfig()
	config.User = username
	config.Passwd = password
	config.Addr = address
	config.DBName = database

	db, err := sql.Open("mysql", config.FormatDSN())
	if err != nil {
		fmt.Printf("mysql-db: error opening database connection. Error: %s\n", err)
		return MysqlDatabaseConnection{}, err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("mysql-db: error pinging database. Error: %s\n", err)
		return MysqlDatabaseConnection{}, err
	}

	err = db.Close()
	if err != nil {
		fmt.Printf("mysql-db: error closing database connection. Error: %s\n", err)
		return MysqlDatabaseConnection{}, err
	}

	return MysqlDatabaseConnection{config: config}, nil
}
