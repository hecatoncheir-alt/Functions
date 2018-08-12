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

func (functions FAASFunctions) ReadCategoriesByName(categoryName, language string) []storage.Category {
	functionPath := fmt.Sprintf(
		"%v/%v/%v", functions.FunctionsGateway, "function", "storage-category-read-by-name")

	body := struct {
		Language        string
		CategoryName    string
		DatabaseGateway string
	}{
		Language:        language,
		CategoryName:    categoryName,
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

	var existCategories []storage.Category

	err = json.Unmarshal(decodedResponse, &existCategories)
	if err != nil {
		FAASLogger.Println(err)
		return nil
	}

	return existCategories
}

func (functions FAASFunctions) ReadCategoryByID(companyName, language string) (storage.Category, error) {
	/// TODO
}
