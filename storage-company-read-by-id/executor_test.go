package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestCompanyCanBeReadByID(t *testing.T) {
	IDOfTestedCompany := "0x12"

	executor := Executor{Store: MockStore{}}

	companyFromStore, err := executor.ReadCompanyByID(IDOfTestedCompany, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if companyFromStore.ID != "0x12" {
		t.Fatalf("Expected id of company: '0x12', actual: %v", companyFromStore.ID)
	}

	if companyFromStore.Name != "Test company" {
		t.Fatalf("Expected name of company: 'Test company', actual: %v", companyFromStore.Name)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "companies":[
			  {
				 "uid":"0x12",
				 "companyName":"Test company",
				 "companyIri":"/",
				 "companyIsActive":true
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestCompanyCanBeReadByNameWithError(t *testing.T) {
	IDOfTestedCompany := "0x12"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadCompanyByID(IDOfTestedCompany, "ru")
	if err != ErrCompanyByIDCanNotBeFound {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStore struct {
	storage.Store
}

func (store ErrorMockStore) Query(request string) (response []byte, err error) {
	return []byte(""), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestCompanyCanBeReadByNameAndItCanBeEmpty(t *testing.T) {
	IDOfTestedCompany := "0x12"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadCompanyByID(IDOfTestedCompany, "ru")
	if err != ErrCompanyDoesNotExist {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "companies":[]
		}
	`

	return []byte(resp), nil
}
