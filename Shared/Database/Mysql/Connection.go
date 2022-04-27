package Dbconnection

import (
	"CitadelCore/Shared/Helpers/Reflection"
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

func (connection MysqlDatabaseConnection) ExecuteQuery(query string, resultType interface{}) (interface{}, error) {
	db, err := sql.Open("mysql", connection.config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("mysql-db: error opening database connection. Error: %s\n", err)
	}
	defer db.Close()

	rows := db.QueryRow(query)

	err = rows.Scan(resultType)
	if err != nil {
		return nil, fmt.Errorf("mysql-db: error reading results. Error: %s\n", err)
	}

	return resultType, nil
}

func (connection MysqlDatabaseConnection) ExecuteQuerySingleResult(query string, resultType interface{}, args ...interface{}) (interface{}, error) {
	db, err := sql.Open("mysql", connection.config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("mysql-db: error opening database connection. Error: %s\n", err)
	}
	defer db.Close()

	rows := db.QueryRow(query, args...)

	fields := Reflection.GetArrayOfFields(resultType)
	err = rows.Scan(fields...)
	if err != nil {
		return nil, fmt.Errorf("mysql-db: error reading results. Error: %s\n", err)
	}

	return Reflection.CreateResultFromFields(fields, resultType), nil
}

func (connection MysqlDatabaseConnection) ExecuteQueryMultipleResults(query string, resultType interface{}, resultSet []interface{}, args ...interface{}) ([]interface{}, error) {
	db, err := sql.Open("mysql", connection.config.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("mysql-db: error opening database connection. Error: %s\n", err)
	}
	defer db.Close()

	rows, err := db.Query(query, args...)

	var index = 0
	for rows.Next() {
		fields := Reflection.GetArrayOfFields(resultType)
		err = rows.Scan(fields...)
		if err != nil {
			return nil, fmt.Errorf("mysql-db: error reading results. Error: %s\n", err)
		}

		resultSet[index] = Reflection.CreateResultFromFields(fields, resultType)
		index++
	}

	return resultSet, nil
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
