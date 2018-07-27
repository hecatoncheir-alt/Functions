package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
	)

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	return []byte(""), nil
}

func TestCompanyCanBeReadByName(t *testing.T) {
	executor := Executor{Store: MockStore{}}

	nameOfTestedCompany := "Test company"

	companiesFromStore, err:= executor.ReadCompaniesByName(nameOfTestedCompany, "ru", "")
	if err != nil {
		t.Error(err)
	}


	if len(companiesFromStore) > 1 {
		t.Fail()
	}

	if companiesFromStore[0].Name != nameOfTestedCompany {
		t.Fail()
	}

	if companiesFromStore[0].ID == "" {
		t.Fail()
	}
}
