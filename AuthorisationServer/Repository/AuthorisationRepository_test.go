package Repository_test

import (
	"CitadelCore/AuthorisationServer/Repository"
	"fmt"
	"testing"
)

func TestGetAccountInformation(t *testing.T) {
	accountname := "TEST"
	repo := Repository.InitializeAuthorisationRepository()

	result := repo.GetAccountInformation(accountname)

	resacc := result.Accountname

	fmt.Println(resacc)
	fmt.Println(result.Salt)
	fmt.Println(result.Verifier)
	if resacc != accountname {
		t.Fatalf("Account name did not return as expected. Expected %s, was %s\n", accountname, resacc)
	}
}

func TestGetRealm(t *testing.T) {
	repo := Repository.InitializeAuthorisationRepository()

	result, _ := repo.GetRealms()

	fmt.Println(result)
	for _, v := range result {
		fmt.Println(v.Name)
	}
}
