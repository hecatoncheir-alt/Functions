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
	// ErrCitiesByNameCanNotBeFound means that the cities can't be found in database
	ErrCitiesByNameCanNotBeFound = errors.New("cities by name can not be found")

	// ErrCitiesByNameNotFound means than the cities does not exist in database
	ErrCitiesByNameNotFound = errors.New("cities by name not found")
)

// ReadCitiesByName is a method for get all nodes by categories name
func (executor *Executor) ReadCitiesByName(cityName, language string) ([]storage.City, error) {
	variables := struct {
		CityName string
		Language string
	}{
		CityName: cityName,
		Language: language}

	queryTemplate, err := template.New("ReadCitiesByName").Parse(`{
				cities(func: eq(cityName@{{.Language}}, "{{.CityName}}")) @filter(eq(cityIsActive, true)) {
					uid
					cityName: cityName@{{.Language}}
					cityIsActive
				}
			}`)

	if err != nil {
		ExecutorLogger.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return nil, err
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		ExecutorLogger.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	type citiesInStorage struct {
		AllCitiesFoundedByName []storage.City `json:"cities"`
	}

	var foundedCities citiesInStorage
	err = json.Unmarshal(response, &foundedCities)
	if err != nil {
		ExecutorLogger.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	if len(foundedCities.AllCitiesFoundedByName) == 0 {
		ExecutorLogger.Println(ErrCitiesByNameNotFound)
		return nil, ErrCitiesByNameNotFound
	}

	return foundedCities.AllCitiesFoundedByName, nil
}
