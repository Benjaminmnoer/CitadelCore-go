package Repository

import (
	model "CitadelCore/AuthorisationServer/Repository/Model"
	connection "CitadelCore/Shared/Database/Mysql"
	"encoding/hex"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	username = "citadelcore"
	password = "citadelcore"
	address  = "127.0.0.1:3307"
	database = "auth"
)

type AuthorisationRepository struct {
	dbconnection connection.MysqlDatabaseConnection
	db           *gorm.DB
}

func (authRepo AuthorisationRepository) GetAccountInformation(accountname string) model.AccountInformation {
	value := model.AccountInformation{}
	authRepo.db.Where("accountname = ?", accountname).First(&value)
	return value
}

func (authRepo AuthorisationRepository) GetRealms() ([]model.Realm, error) {
	realms := []model.Realm{}
	authRepo.db.Find(&realms)
	return realms, nil
}

func InitializeAuthorisationRepository() AuthorisationRepository {
	db, err := gorm.Open(mysql.Open("citadelcore:citadelcore@tcp(127.0.0.1:3307)/auth?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.AccountInformation{}, &model.Realm{})

	salt, _ := hex.DecodeString("2C9B1534B3E0D354EB682BF203E76D67BE621399F4DBFE8054C84E3D2977398E")
	for i, j := 0, len(salt)-1; i < j; i, j = i+1, j-1 {
		salt[i], salt[j] = salt[j], salt[i]
	}

	verifier, _ := hex.DecodeString("E475392E52BBF123D14780189E7AF2D1C051FE6A0476360FB012C93DB715EFAB")
	for i, j := 0, len(verifier)-1; i < j; i, j = i+1, j-1 {
		verifier[i], verifier[j] = verifier[j], verifier[i]
	}

	db.FirstOrCreate(&model.AccountInformation{Id: 0, Accountname: "TEST", Salt: salt, Verifier: verifier})
	db.FirstOrCreate(&model.Realm{Id: 0, Name: "Trintiy", Address: "127.0.0.1", Port: 8085, Icon: 0, Flag: 2, Timezone: 1, AllowedSecurityLevel: 0, Population: 0, Gamebuild: 12340})

	if err != nil {
		panic(err)
	}

	return AuthorisationRepository{db: db}
}
