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

func (functions FAASFunctions) ReadPriceByID(priceID, language string) storage.Price {
	functionPath := fmt.Sprintf(
		"%v/%v/%v", functions.FunctionsGateway, "function", "storage-price-read-by-id")

	body := struct {
		Language        string
		PriceID         string
		DatabaseGateway string
	}{
		Language:        language,
		PriceID:         priceID,
		DatabaseGateway: functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Price{}
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		FAASLogger.Println(err)
		return storage.Price{}
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Price{}
	}

	var existPrice storage.Price

	err = json.Unmarshal(decodedResponse, &existPrice)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Price{}
	}

	return existPrice
}
