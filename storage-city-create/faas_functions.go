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

var logger = log.New(os.Stdout, "FAASFunctions: ", log.Lshortfile)

type FAASFunctions struct{ FunctionsGateway, DatabaseGateway string }

func (functions FAASFunctions) ReadCitiesByName(cityName, language string) []storage.City {
	functionPath := fmt.Sprintf(
		"%v/%v", functions.FunctionsGateway, "storage-city-read-by-name")

	body := struct {
		Language, CityName, DatabaseGateway string
	}{
		Language:        language,
		CityName:        cityName,
		DatabaseGateway: functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		logger.Println(err)
		return nil
	}

	responseWithEncodedBody, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		logger.Println(err)
		return nil
	}

	defer responseWithEncodedBody.Body.Close()

	decodedResponse, err := ioutil.ReadAll(responseWithEncodedBody.Body)
	if err != nil {
		logger.Println(err)
		return nil
	}

	response := Response{}

	err = json.Unmarshal([]byte(decodedResponse), &response)
	if err != nil {
		logger.Println(err)
		return nil
	}

	if response.Error != "" {
		logger.Println(err)
		return nil
	}

	var existCities []storage.City

	err = json.Unmarshal([]byte(response.Data), &existCities)
	if err != nil {
		logger.Println(err)
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
		logger.Println(err)
		return storage.City{}
	}

	responseWithEncodedBody, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		logger.Println(err)
		return storage.City{}
	}

	defer responseWithEncodedBody.Body.Close()

	decodedResponse, err := ioutil.ReadAll(responseWithEncodedBody.Body)
	if err != nil {
		logger.Println(err)
		return storage.City{}
	}

	response := Response{}

	err = json.Unmarshal(decodedResponse, &response)
	if err != nil {
		logger.Println(err)
		return storage.City{}
	}

	if response.Error != "" {
		logger.Println(err)
		return storage.City{}
	}

	var existCity storage.City

	err = json.Unmarshal([]byte(response.Data), &existCity)
	if err != nil {
		logger.Println(err)
		return storage.City{}
	}

	return existCity
}
