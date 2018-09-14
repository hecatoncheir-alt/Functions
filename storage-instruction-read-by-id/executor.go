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
	// ErrInstructionCanNotBeWithoutID means that instruction can't be without id
	ErrInstructionCanNotBeWithoutID = errors.New("instruction can not be without id")

	// ErrInstructionByIDCanNotBeFound means that the instruction can't be found in database
	ErrInstructionByIDCanNotBeFound = errors.New("instruction by id can not be found")

	// ErrInstructionDoesNotExist means than the instruction does not exist in database
	ErrInstructionDoesNotExist = errors.New("instruction does not exist")
)

// ReadInstructionByID is a method for get all nodes of instructions by ID
func (executor *Executor) ReadInstructionByID(instructionID, language string) (storage.Instruction, error) {

	if instructionID == "" {
		ExecutorLogger.Printf("Instruction can't be without ID")
		return storage.Instruction{}, ErrInstructionCanNotBeWithoutID
	}

	variables := struct {
		InstructionID string
		Language      string
	}{
		InstructionID: instructionID,
		Language:      language}

	queryTemplate, err := template.New("ReadInstructionByID").Parse(`{
				instructions(func: uid("{{.InstructionID}}")) @filter(eq(instructionLanguage, {{.Language}})) {
					uid
					instructionLanguage
					instructionIsActive
					has_page {
						uid
						path
						pageInPaginationSelector
						pageParamPath
						cityParamPath
						itemSelector
						nameOfItemSelector
						priceOfItemSelector
					}
					has_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@{{.Language}}
						cityIsActive
					}
					has_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIsActive
					}
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
					}
				}
			}`)

	instruction := storage.Instruction{ID: instructionID}

	if err != nil {
		ExecutorLogger.Println(err)
		return instruction, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return instruction, err
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		ExecutorLogger.Println(err)
		return instruction, ErrInstructionByIDCanNotBeFound
	}

	type InstructionsInStorage struct {
		Instructions []storage.Instruction `json:"instructions"`
	}

	var foundedInstructions InstructionsInStorage

	err = json.Unmarshal(response, &foundedInstructions)
	if err != nil {
		ExecutorLogger.Println(err)
		return instruction, ErrInstructionByIDCanNotBeFound
	}

	if len(foundedInstructions.Instructions) == 0 {
		return instruction, ErrInstructionDoesNotExist
	}

	return foundedInstructions.Instructions[0], nil
}
