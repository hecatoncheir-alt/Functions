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

func (functions FAASFunctions) ReadInstructionByID(instructionID, language string) storage.Instruction {
	functionPath := fmt.Sprintf(
		"%v/%v/%v", functions.FunctionsGateway, "function", "storage-instruction-read-by-id")

	body := struct {
		Language        string
		InstructionID   string
		DatabaseGateway string
	}{
		Language:        language,
		InstructionID:   instructionID,
		DatabaseGateway: functions.DatabaseGateway}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Instruction{}
	}

	response, err := http.Post(functionPath, "application/json", bytes.NewBuffer(encodedBody))
	if err != nil {
		FAASLogger.Println(err)
		return storage.Instruction{}
	}

	defer response.Body.Close()

	decodedResponse, err := ioutil.ReadAll(response.Body)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Instruction{}
	}

	var existInstruction storage.Instruction

	err = json.Unmarshal(decodedResponse, &existInstruction)
	if err != nil {
		FAASLogger.Println(err)
		return storage.Instruction{}
	}

	return existInstruction
}
