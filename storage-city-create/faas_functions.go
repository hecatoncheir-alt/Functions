package function

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var FAASLogger = log.New(os.Stdout, "FAASFunctions: ", log.Lshortfile)

type FAASFunctions struct {
	FunctionsGateway string
	DatabaseGateway  string
}

func (functions FAASFunctions) ReadCitiesByName(cityName, language string) []storage.City {
	functionPath := fmt.Sprintf(
		"%v/%v", functions.FunctionsGateway, "storage-city-read-by-name")

	body := struct {
		Language        string
		CityName        string
		DatabaseGateway string
	}{
		Language:        language,
		CityName:        cityName,
		DatabaseGateway: functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
		return nil
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		FAASLogger.Println(err)
		return nil
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		FAASLogger.Println(err)
		return nil
	}

	var existCities []storage.City

	err = json.Unmarshal(decodedResponse, &existCities)
	if err != nil {
		FAASLogger.Println(err)
		return nil
	}

	return existCities
}

func (functions FAASFunctions) ReadCityByID(cityID, language string) storage.City {
	functionPath := fmt.Sprintf(
		"%v/%v/%v", functions.FunctionsGateway, "function", "storage-city-read-by-id")

	body := struct {
		Language        string
		CityID          string
		DatabaseGateway string
	}{
		Language:        language,
		CityID:          cityID,
		DatabaseGateway: functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.City{}
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		FAASLogger.Println(err)
		return storage.City{}
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.City{}
	}

	var existCity storage.City

	err = json.Unmarshal(decodedResponse, &existCity)
	if err != nil {
		FAASLogger.Println(err)
		return storage.City{}
	}

	return existCity
}
