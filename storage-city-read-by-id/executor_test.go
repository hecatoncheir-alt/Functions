package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestCityCanBeReadByID(t *testing.T) {
	IDOfTestedCity := "0x12"

	executor := Executor{Store: MockStore{}}

	cityFromStore, err := executor.ReadCityByID(IDOfTestedCity, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if cityFromStore.ID != "0x12" {
		t.Fatalf("Expected id of city: '0x12', actual: %v", cityFromStore.ID)
	}

	if cityFromStore.Name != "Test city" {
		t.Fatalf("Expected name of city: 'Test city', actual: %v", cityFromStore.Name)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "cities":[
			  {
				 "uid":"0x12",
				 "cityName":"Test city",
				 "cityIsActive":true
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestCityCanBeReadByNameWithError(t *testing.T) {
	IDOfTestedCity := "0x12"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadCityByID(IDOfTestedCity, "ru")
	if err != ErrCityByIDCanNotBeFound {
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

func TestCityCanBeReadByNameAndItCanBeEmpty(t *testing.T) {
	IDOfTestedCity := "0x12"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadCityByID(IDOfTestedCity, "ru")
	if err != ErrCityDoesNotExist {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "cities":[]
		}
	`

	return []byte(resp), nil
}
