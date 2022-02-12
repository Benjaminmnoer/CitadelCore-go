package Repository_test

import (
	"CitadelCore/AuthorisationServer/Repository"
	"testing"
)

func TestGetAccountInformation(t *testing.T) {
	accountname := "TEST"
	repo := Repository.InitializeAuthorisationRepository()

	result := repo.GetAccountInformation(accountname)

	resacc := result.Accountname

	if resacc != accountname {
		t.Fatalf("Account name did not return as expected. Expected %s, was %s\n", accountname, resacc)
	}
}
