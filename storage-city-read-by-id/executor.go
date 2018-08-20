package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
	"text/template"
)

type Storage interface {
	Query(string) ([]byte, error)
}

type Executor struct {
	Store Storage
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCityCanNotBeWithoutID means that city can't be without id
	ErrCityCanNotBeWithoutID = errors.New("city can not be without id")

	// ErrCityByIDCanNotBeFound means that the city can't be found in database
	ErrCityByIDCanNotBeFound = errors.New("city by id can not be found")

	// ErrCityDoesNotExist means than the city does not exist in database
	ErrCityDoesNotExist = errors.New("city does not exist")
)

// ReadCityByID is a method for get all nodes of city by ID
func (executor *Executor) ReadCityByID(cityID, language string) (storage.City, error) {
	city := storage.City{}

	if cityID == "" {
		ExecutorLogger.Printf("City can't be without ID")
		return city, ErrCityCanNotBeWithoutID
	}

	variables := struct {
		CityID   string
		Language string
	}{
		CityID:   cityID,
		Language: language}

	queryTemplate, err := template.New("ReadCityByID").Parse(`{
				cities(func: uid("{{.CityID}}")) @filter(has(cityName)) {
					uid
					cityName: cityName@{{.Language}}
					cityIsActive
				}
			}`)

	city = storage.City{ID: cityID}
	if err != nil {
		ExecutorLogger.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return city, err
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		ExecutorLogger.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	type citiesInStore struct {
		Cities []storage.City `json:"cities"`
	}

	var foundedCities citiesInStore

	err = json.Unmarshal(response, &foundedCities)
	if err != nil {
		ExecutorLogger.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	if len(foundedCities.Cities) == 0 {
		return city, ErrCityDoesNotExist
	}

	return foundedCities.Cities[0], nil
}
