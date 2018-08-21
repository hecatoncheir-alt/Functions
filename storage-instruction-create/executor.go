package function

import (
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
)

type Storage interface {
	Mutate([]byte) (string, error)
	SetNQuads(string, string, string) error
}

type Functions interface {
	ReadInstructionByID(string, string) storage.Instruction
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrInstructionCanNotBeCreated means that the instruction can't be added to database
	ErrInstructionCanNotBeCreated = errors.New("instruction can't be created")
)

// CreateInstruction make category and save it to storage
func (executor *Executor) CreateInstruction(instruction storage.Instruction, language string) (storage.Instruction, error) {

	instruction.IsActive = true

	encodedInstruction, err := json.Marshal(instruction)
	if err != nil {
		return instruction, ErrInstructionCanNotBeCreated
	}

	uidOfCreatedInstruction, err := executor.Store.Mutate(encodedInstruction)
	if err != nil {
		return instruction, ErrInstructionCanNotBeCreated
	}

	createdCity := executor.Functions.ReadInstructionByID(uidOfCreatedInstruction, language)

	return createdCity, nil
}
