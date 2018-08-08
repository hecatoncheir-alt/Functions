package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCompanyCanBeCreated(t *testing.T) {
	companyForCreate := storage.Company{ID: "0x12", Name: "Test company", IsActive: true}

	executor := Executor{
		Functions: EmptyCompaniesFAASFunctions{FunctionsGateway: ""},
		Store:     MockStorage{DatabaseGateway: ""}}

	createdCompany, err := executor.CreateCompany(companyForCreate, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCompany.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", companyForCreate.ID, createdCompany.ID)
	}

	if createdCompany.Name != "Test company" {
		t.Errorf("Expect: %v, but got: %v", companyForCreate.Name, createdCompany.Name)
	}
}

/// Mock FAAS functions
type EmptyCompaniesFAASFunctions struct {
	FunctionsGateway string
}

func (functions EmptyCompaniesFAASFunctions) CompaniesReadByName(companyName, language string) []storage.Company {
	return []storage.Company{}
}

func (functions EmptyCompaniesFAASFunctions) ReadCompanyByID(companyID, language string) (storage.Company, error) {
	return storage.Company{ID: companyID, Name: "Test company"}, nil
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) Mutate(setJson []byte) (uid string, err error) {
	return "0x12", nil
}

func (store MockStorage) SetNQuads(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestCompanyCanNotBeCreated(t *testing.T) {
	// TODO
}

// --------------------------------------------------------------------------------------------------------
func TestCreatingCompanyCanBeExists(t *testing.T) {

	companyForCreate := storage.Company{ID: "0x12", Name: "Test company", IsActive: true}

	executor := Executor{Functions: NotEmptyCompaniesFAASFunctions{}}

	existFirstTestCompany, err := executor.CreateCompany(companyForCreate, "ru")
	if err != ErrCompanyAlreadyExist {
		t.Fatalf(err.Error())
	}

	if existFirstTestCompany.Name != "First test company name" {
		t.Errorf("Expect: %v, but got: %v", "First test company name", existFirstTestCompany.Name)
	}

}

type NotEmptyCompaniesFAASFunctions struct{}

func (functions NotEmptyCompaniesFAASFunctions) CompaniesReadByName(companyName, language string) []storage.Company {
	return []storage.Company{
		{ID: "0x12", Name: "First test company name", IsActive: true},
		{ID: "0x13", Name: "Second test company name", IsActive: true}}
}

func (functions NotEmptyCompaniesFAASFunctions) ReadCompanyByID(companyID, language string) (storage.Company, error) {
	return storage.Company{ID: "0x13", Name: "Second test company name"}, nil
}

// ---------------------------------------------------------------------------------------------------------------
