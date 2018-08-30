package function

import (
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
)

type Storage interface {
	CreateJSON([]byte) (string, error)
	AddLanguage(string, string, string) error
}

type Functions interface {
	ReadCitiesByName(string, string) []storage.City
	ReadCityByID(string, string) storage.City
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCityCanNotBeCreated means that the city can't be added to database
	ErrCityCanNotBeCreated = errors.New("city can't be created")

	// ErrCityAlreadyExist means that the city is in the database already
	ErrCityAlreadyExist = errors.New("city already exist")
)

// CreateCity make category and save it to storage
func (executor *Executor) CreateCity(city storage.City, language string) (storage.City, error) {

	existsCities := executor.Functions.ReadCitiesByName(city.Name, language)

	if len(existsCities) > 0 {
		ExecutorLogger.Printf("City with name: %v exist: %v", city.Name, existsCities[0])
		return existsCities[0], ErrCityAlreadyExist
	}

	city.IsActive = true

	encodedCity, err := json.Marshal(city)
	if err != nil {
		return city, ErrCityCanNotBeCreated
	}

	uidOfCreatedCity, err := executor.Store.CreateJSON(encodedCity)
	if err != nil {
		return city, ErrCityCanNotBeCreated
	}

	err = executor.Store.AddLanguage(uidOfCreatedCity, "cityName", "\""+city.Name+"\""+"@"+language)
	if err != nil {
		return city, ErrCityCanNotBeCreated
	}

	createdCity := executor.Functions.ReadCityByID(uidOfCreatedCity, language)

	return createdCity, nil
}
