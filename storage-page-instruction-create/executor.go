package function

import (
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
)

type Storage interface {
	CreateJSON([]byte) (string, error)
}

type Functions interface {
	ReadPageInstructionByID(string, string) storage.PageInstruction
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrPageInstructionCanNotBeCreated means that the page instruction can't be added to database
	ErrPageInstructionCanNotBeCreated = errors.New("page instruction can't be created")
)

// CreatePageInstruction make PageInstruction and save it to storage
func (executor *Executor) CreatePageInstruction(pageInstruction storage.PageInstruction, language string) (storage.PageInstruction, error) {

	encodedPageInstruction, err := json.Marshal(pageInstruction)
	if err != nil {
		return pageInstruction, ErrPageInstructionCanNotBeCreated
	}

	uidOfCreatedPageInstruction, err := executor.Store.CreateJSON(encodedPageInstruction)
	if err != nil {
		return pageInstruction, ErrPageInstructionCanNotBeCreated
	}

	createdPageInstruction := executor.Functions.ReadPageInstructionByID(uidOfCreatedPageInstruction, language)

	return createdPageInstruction, nil
}
