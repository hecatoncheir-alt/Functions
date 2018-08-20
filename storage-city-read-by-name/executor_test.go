package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

func TestCityCanBeReadByName(t *testing.T) {
	nameOfTestedCity := "Test city"

	executor := Executor{Store: MockStore{}}

	citiesFromStore, err := executor.ReadCitiesByName(nameOfTestedCity, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(citiesFromStore) < 1 {
		t.Fatalf("Expected 1 city, actual: %v", len(citiesFromStore))
	}

	if citiesFromStore[0].ID != "0x12" {
		t.Fatalf("Expected id of city: '0x12', actual: %v", citiesFromStore[0].ID)
	}

	if citiesFromStore[0].Name != nameOfTestedCity {
		t.Fatalf("Expected name of city: 'Test city', actual: %v", citiesFromStore[0].Name)
	}
}

func TestCityCanBeReadByNameWithError(t *testing.T) {
	nameOfTestedCity := "Test city"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadCitiesByName(nameOfTestedCity, "ru")
	if err != ErrCitiesByNameCanNotBeFound {
		t.Fatalf(err.Error())
	}
}

func TestCityCanBeReadByNameAndItCanBeEmpty(t *testing.T) {
	nameOfTestedCity := "Test city"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadCitiesByName(nameOfTestedCity, "ru")
	if err != ErrCitiesByNameNotFound {
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
		   "cities":[]
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
		   "cities":[  
			  {  
				 "uid":"0x12",
				 "cityName":"Test city",
				 "cityIsActive":true
			  },
			  {  
				 "uid":"0x13",
				 "cityName":"Other test city",
				 "cityIsActive":true
			  }
		   ]
		}
	`

	return []byte(resp), nil
}
