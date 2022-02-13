package Repository

import (
	model "CitadelCore/AuthorisationServer/Model"
	connection "CitadelCore/Shared/Database/Mysql"
	"fmt"
)

var (
	username = "trinity"
	password = "trinity"
	address  = "127.0.0.1"
	database = "auth"
)

const (
	_ACCOUNTINFO_QUERY = "SELECT username, salt, verifier FROM account WHERE username = ?;"
)

type AuthorisationRepository struct {
	dbconnection connection.MysqlDatabaseConnection
}

func (authRepo AuthorisationRepository) GetAccountInformation(accountname string) model.AccountInformation {
	fmt.Printf("Querying for name: %s\n", accountname)
	accountinfo := model.AccountInformation{}
	_, err := authRepo.dbconnection.ExecuteQuerySingleResult(_ACCOUNTINFO_QUERY, &accountinfo, accountname)
	if err != nil {
		panic(err)
	}

	return accountinfo
}

func InitializeAuthorisationRepository() AuthorisationRepository {
	conn, err := connection.CreateDatabaseConnection(username, password, address, database)

	if err != nil {
		panic(err)
	}

	return AuthorisationRepository{dbconnection: conn}
}
