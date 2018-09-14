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

func (functions FAASFunctions) CompaniesReadByName(companyName, language string) []storage.Company {
	functionPath := fmt.Sprintf(
		"%v/%v", functions.FunctionsGateway, "storage-company-read-by-name")

	body := struct {
		Language        string
		CompanyName     string
		DatabaseGateway string
	}{
		Language:        language,
		CompanyName:     companyName,
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

	var existCompanies []storage.Company

	err = json.Unmarshal(decodedResponse, &existCompanies)
	if err != nil {
		FAASLogger.Println(err)
		return nil
	}

	return existCompanies
}

func (functions FAASFunctions) ReadCompanyByID(companyID, language string) storage.Company {
	functionPath := fmt.Sprintf(
		"%v/%v/%v", functions.FunctionsGateway, "function", "storage-company-read-by-id")

	body := struct {
		Language        string
		CompanyID       string
		DatabaseGateway string
	}{
		Language:        language,
		CompanyID:       companyID,
		DatabaseGateway: functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Company{}
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		FAASLogger.Println(err)
		return storage.Company{}
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Company{}
	}

	var existCompany storage.Company

	err = json.Unmarshal(decodedResponse, &existCompany)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Company{}
	}

	return existCompany
}
