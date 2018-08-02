package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

type EmptyMockStore struct {
	Store Storage
	Functions Functions
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		[]
	`

	return []byte(resp), nil
}

func TestCompanyCanNotBeCreated(t *testing.T) {

}

type EmptyCompaniesFAASFunctions struct {}

func (functions EmptyCompaniesFAASFunctions) CompaniesReadByName(companyName, language, DatabaseGateway string) []storage.Company {
	return []storage.Company{
		{ID:"0x12", Name:"First test company name", IsActive: true},
		{ID:"0x13", Name:"Second test company name", IsActive: true}}
}


func TestCompanyCanBeExists(t *testing.T) {

	companyForCreate := storage.Company{ID: "0x12", Name: "Test company", IsActive: true}

	executor := Executor{Functions: EmptyCompaniesFAASFunctions{}}

	existFirstTestCompany, err := executor.CreateCompany(companyForCreate, "ru", "//")
	if err != ErrCompanyAlreadyExist {
		t.Fatalf(err.Error())
	}

	if existFirstTestCompany.Name  != "First test company name" {
		t.Errorf("Expect: %v, but got: %v", "First test company name", existFirstTestCompany.Name)
	}

}
