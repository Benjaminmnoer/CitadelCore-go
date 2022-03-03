package Repository

import (
	model "CitadelCore/AuthorisationServer/Model"
	connection "CitadelCore/Shared/Database/Mysql"
	"fmt"
)

var (
	username = "citadelcore"
	password = "citadelcore"
	address  = "127.0.0.1:3306"
	database = "auth"
)

const (
	_ACCOUNTINFO_QUERY = "SELECT username, salt, verifier FROM account WHERE username = ?;"
	_REALMLIST_QUERY   = "SELECT id, name, address, localAddress, localSubnetMask, port, icon, flag, timezone, allowedSecurityLevel, population, gamebuild FROM realmlist WHERE flag <> 3 ORDER BY name;"
)

type AuthorisationRepository struct {
	dbconnection connection.MysqlDatabaseConnection
}

func (authRepo AuthorisationRepository) GetAccountInformation(accountname string) model.AccountInformation {
	query := _ACCOUNTINFO_QUERY // strings.Replace(_ACCOUNTINFO_QUERY, "?", accountname, 1)
	fmt.Println(query)
	accountinfo := model.AccountInformation{}
	_, err := authRepo.dbconnection.ExecuteQuerySingleResult(query, &accountinfo, accountname)
	if err != nil {
		panic(err)
	}
	for i, j := 0, len(accountinfo.Verifier)-1; i < j; i, j = i+1, j-1 {
		accountinfo.Verifier[i], accountinfo.Verifier[j] = accountinfo.Verifier[j], accountinfo.Verifier[i]
	}
	for i, j := 0, len(accountinfo.Salt)-1; i < j; i, j = i+1, j-1 {
		accountinfo.Salt[i], accountinfo.Salt[j] = accountinfo.Salt[j], accountinfo.Salt[i]
	}

	return accountinfo
}

func (authRepo AuthorisationRepository) GetRealms() {

}

func InitializeAuthorisationRepository() AuthorisationRepository {
	conn, err := connection.CreateDatabaseConnection(username, password, address, database)

	if err != nil {
		panic(err)
	}

	return AuthorisationRepository{dbconnection: conn}
}
