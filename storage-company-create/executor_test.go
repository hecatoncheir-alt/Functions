package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

type EmptyCompaniesFAASFunctions struct{}

func (functions EmptyCompaniesFAASFunctions) CompaniesReadByName(
	companyName, language, DatabaseGateway string) []storage.Company {
	return []storage.Company{}
}

func TestCompanyCanBeCreated(t *testing.T) {
	companyForCreate := storage.Company{ID: "0x12", Name: "Test company", IsActive: true}

	executor := Executor{Functions: EmptyCompaniesFAASFunctions{}}

	createdCompany, err := executor.CreateCompany(companyForCreate, "ru", "//")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCompany.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", companyForCreate.ID, createdCompany.Name)
	}

	if createdCompany.Name != "Test company" {
		t.Errorf("Expect: %v, but got: %v", companyForCreate.Name, createdCompany.Name)
	}
}

func TestCompanyCanNotBeCreated(t *testing.T) {
	// TODO
}

type NotEmptyCompaniesFAASFunctions struct{}

func (functions NotEmptyCompaniesFAASFunctions) CompaniesReadByName(companyName, language, DatabaseGateway string) []storage.Company {
	return []storage.Company{
		{ID: "0x12", Name: "First test company name", IsActive: true},
		{ID: "0x13", Name: "Second test company name", IsActive: true}}
}

func TestCreatingCompanyCanBeExists(t *testing.T) {

	companyForCreate := storage.Company{ID: "0x12", Name: "Test company", IsActive: true}

	executor := Executor{Functions: NotEmptyCompaniesFAASFunctions{}}

	existFirstTestCompany, err := executor.CreateCompany(companyForCreate, "ru", "//")
	if err != ErrCompanyAlreadyExist {
		t.Fatalf(err.Error())
	}

	if existFirstTestCompany.Name != "First test company name" {
		t.Errorf("Expect: %v, but got: %v", "First test company name", existFirstTestCompany.Name)
	}

}
