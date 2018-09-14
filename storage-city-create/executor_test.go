package function

import (
	"errors"
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCityCanBeCreated(t *testing.T) {
	cityForCreate := storage.City{ID: "0x12", Name: "Test city", IsActive: true}

	executor := Executor{
		Functions: EmptyCitiesFAASFunctions{FunctionsGateway: ""},
		Store:     MockStorage{DatabaseGateway: ""}}

	createdCity, err := executor.CreateCity(cityForCreate, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCity.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", cityForCreate.ID, createdCity.ID)
	}

	if createdCity.Name != "Test city" {
		t.Errorf("Expect: %v, but got: %v", cityForCreate.Name, createdCity.Name)
	}
}

/// Mock FAAS functions
type EmptyCitiesFAASFunctions struct {
	FunctionsGateway string
}

func (functions EmptyCitiesFAASFunctions) ReadCitiesByName(companyName, language string) []storage.City {
	return []storage.City{}
}

func (functions EmptyCitiesFAASFunctions) ReadCityByID(cityID, language string) storage.City {
	return storage.City{ID: cityID, Name: "Test city"}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "0x12", nil
}

func (store MockStorage) AddLanguage(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestCityCanNotBeCreated(t *testing.T) {
	cityForCreate := storage.City{ID: "0x12", Name: "Test city", IsActive: true}

	executor := Executor{
		Functions: EmptyCitiesFAASFunctions{FunctionsGateway: ""},
		Store:     ErrorMockStorage{DatabaseGateway: ""}}

	_, err := executor.CreateCity(cityForCreate, "ru")
	if err != ErrCityCanNotBeCreated {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStorage struct {
	DatabaseGateway string
}

func (store ErrorMockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "", errors.New("")
}

func (store ErrorMockStorage) AddLanguage(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------
func TestCreatingCityCanBeExists(t *testing.T) {

	cityForCreate := storage.City{ID: "0x12", Name: "Test city", IsActive: true}

	executor := Executor{Functions: NotEmptyCitiesFAASFunctions{}}

	existFirstTestCity, err := executor.CreateCity(cityForCreate, "ru")
	if err != ErrCityAlreadyExist {
		t.Fatalf(err.Error())
	}

	if existFirstTestCity.Name != "First test city name" {
		t.Errorf("Expect: %v, but got: %v", "First test city name", existFirstTestCity.Name)
	}

}

type NotEmptyCitiesFAASFunctions struct{}

func (functions NotEmptyCitiesFAASFunctions) ReadCitiesByName(cityName, language string) []storage.City {
	return []storage.City{
		{ID: "0x12", Name: "First test city name", IsActive: true},
		{ID: "0x13", Name: "Second test city name", IsActive: true}}
}

func (functions NotEmptyCitiesFAASFunctions) ReadCityByID(cityID, language string) storage.City {
	return storage.City{ID: "0x13", Name: "Second test city name"}
}
