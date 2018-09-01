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

func (functions FAASFunctions) ReadProductsByName(productName, language string) []storage.Product {
	functionPath := fmt.Sprintf(
		"%v/%v", functions.FunctionsGateway, "storage-product-read-by-name")

	body := struct {
		Language        string
		ProductName     string
		DatabaseGateway string
	}{
		Language:        language,
		ProductName:     productName,
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

	var existProducts []storage.Product

	err = json.Unmarshal(decodedResponse, &existProducts)
	if err != nil {
		FAASLogger.Println(err)
		return nil
	}

	return existProducts
}

func (functions FAASFunctions) ReadProductByID(productID, language string) storage.Product {
	functionPath := fmt.Sprintf(
		"%v/%v/%v", functions.FunctionsGateway, "function", "storage-product-read-by-id")

	body := struct {
		Language        string
		ProductID       string
		DatabaseGateway string
	}{
		Language:        language,
		ProductID:       productID,
		DatabaseGateway: functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Product{}
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		FAASLogger.Println(err)
		return storage.Product{}
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Product{}
	}

	var existProduct storage.Product

	err = json.Unmarshal(decodedResponse, &existProduct)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Product{}
	}

	return existProduct
}
