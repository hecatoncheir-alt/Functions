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

func (functions FAASFunctions) ReadPageInstructionByID(pageInstructionID, language string) storage.PageInstruction {
	functionPath := fmt.Sprintf(
		"%v/%v", functions.FunctionsGateway, "storage-page-instruction-read-by-id")

	body := struct {
		Language          string
		PageInstructionID string
		DatabaseGateway   string
	}{
		Language:          language,
		PageInstructionID: pageInstructionID,
		DatabaseGateway:   functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.PageInstruction{ID: pageInstructionID}
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		FAASLogger.Println(err)
		return storage.PageInstruction{ID: pageInstructionID}
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.PageInstruction{ID: pageInstructionID}
	}

	var existPageInstruction storage.PageInstruction

	err = json.Unmarshal(decodedResponse, &existPageInstruction)
	if err != nil {
		FAASLogger.Println(err)
		return storage.PageInstruction{ID: pageInstructionID}
	}

	return existPageInstruction
}
