package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

func TestCompanyCanBeReadByName(t *testing.T) {
	nameOfTestedCompany := "Test company"

	executor := Executor{Store: MockStore{}}

	companiesFromStore, err := executor.ReadCompaniesByName(nameOfTestedCompany, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(companiesFromStore) < 1 {
		t.Fatalf("Expected 1 company, actual: %v", len(companiesFromStore))
	}

	if companiesFromStore[0].ID != "0x12" {
		t.Fatalf("Expected id of company: '0x12', actual: %v", companiesFromStore[0].ID)
	}

	if companiesFromStore[0].Name != nameOfTestedCompany {
		t.Fatalf("Expected name of company: 'Test company', actual: %v", companiesFromStore[0].Name)
	}
}

func TestCompanyCanBeReadByNameWithError(t *testing.T) {
	nameOfTestedCompany := "Test company"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadCompaniesByName(nameOfTestedCompany, "ru")
	if err != ErrCompaniesByNameCanNotBeFound {
		t.Fatalf(err.Error())
	}
}

func TestCompanyCanBeReadByNameAndItCanBeEmpty(t *testing.T) {
	nameOfTestedCompany := "Test company"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadCompaniesByName(nameOfTestedCompany, "ru")
	if err != ErrCompaniesByNameNotFound {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStore struct {
	storage.Store
}

func (store ErrorMockStore) Query(request string) (response []byte, err error) {
	return []byte(""), nil
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
			  },
			  {  
				 "uid":"0x13",
				 "companyName":"Other test company",
				 "companyIri":"/",
				 "companyIsActive":true
			  }
		   ]
		}
	`

	return []byte(resp), nil
}
