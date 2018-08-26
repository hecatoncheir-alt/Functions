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
	// ErrPageInstructionCanNotBeWithoutID means that page instruction can't be without id
	ErrPageInstructionCanNotBeWithoutID = errors.New("page instruction can not be without id")

	// ErrPageInstructionByIDCanNotBeFound means that the instruction can't be found in database
	ErrPageInstructionByIDCanNotBeFound = errors.New("page instruction by id can not be found")

	// ErrPageInstructionDoesNotExist means than the page instruction does not exist in database
	ErrPageInstructionDoesNotExist = errors.New("page instruction does not exist")
)

// ReadPageInstructionByID is a method for get all nodes of instructions by ID
func (executor *Executor) ReadPageInstructionByID(pageInstructionID, language string) (storage.PageInstruction, error) {

	if pageInstructionID == "" {
		ExecutorLogger.Printf("Page pageInstruction can't be without ID")
		return storage.PageInstruction{}, ErrPageInstructionCanNotBeWithoutID
	}

	variables := struct {
		PageInstructionID string
		Language          string
	}{
		PageInstructionID: pageInstructionID,
		Language:          language}

	queryTemplate, err := template.New("ReadPageInstructionByID").Parse(`{
				pageInstructions(func: uid("{{.PageInstructionID}}")) @filter(has(path)) {
					uid
					path
					pageInPaginationSelector
					pageParamPath
					cityParamPath
					itemSelector
					nameOfItemSelector
					priceOfItemSelector
				}
			}`)

	pageInstruction := storage.PageInstruction{ID: pageInstructionID}

	if err != nil {
		ExecutorLogger.Println(err)
		return pageInstruction, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return pageInstruction, err
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		ExecutorLogger.Println(err)
		return pageInstruction, ErrPageInstructionByIDCanNotBeFound
	}

	type PageInstructionsInStorage struct {
		PageInstructions []storage.PageInstruction `json:"pageInstructions"`
	}

	var foundedPageInstructions PageInstructionsInStorage

	err = json.Unmarshal(response, &foundedPageInstructions)
	if err != nil {
		ExecutorLogger.Println(err)
		return pageInstruction, ErrPageInstructionByIDCanNotBeFound
	}

	if len(foundedPageInstructions.PageInstructions) == 0 {
		return pageInstruction, ErrPageInstructionDoesNotExist
	}

	return foundedPageInstructions.PageInstructions[0], nil
}
