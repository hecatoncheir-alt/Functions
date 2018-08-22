package function

import (
	"errors"
	"log"
	"os"
)

type Storage interface {
	SetNQuads(string, string, string) error
}

type Executor struct {
	Store Storage
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrPageInstructionCanNotBeAddedToInstruction means that the PageInstruction can't be added to instruction
	ErrPageInstructionCanNotBeAddedToInstruction = errors.New("page instruction can not be added to instruction")
)

// AddPageInstructionToInstruction method for set quad of predicate about PageInstruction and Instruction
func (executor *Executor) AddPageInstructionToInstruction(instructionID, pageInstructionID string) error {
	err := executor.Store.SetNQuads(instructionID, "has_page", pageInstructionID)
	if err != nil {
		ExecutorLogger.Printf("PageInstruction with ID: %v can not be added to instruction with ID: %v", pageInstructionID, instructionID)
		return ErrPageInstructionCanNotBeAddedToInstruction
	}

	return nil
}
