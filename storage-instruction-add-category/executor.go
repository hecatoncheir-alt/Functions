package function

import (
	"errors"
	"log"
	"os"
)

type Storage interface {
	AddEntityToOtherEntity(string, string, string) error
}

type Executor struct {
	Store Storage
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCategoryCanNotBeAddedToInstruction means that the Category can't be added to instruction
	ErrCategoryCanNotBeAddedToInstruction = errors.New("category can not be added to instruction")
)

// AddCategoryToInstruction method for set quad of predicate about Category and Instruction
func (executor *Executor) AddCategoryToInstruction(instructionID, categoryID string) error {
	err := executor.Store.AddEntityToOtherEntity(instructionID, "has_category", categoryID)
	if err != nil {
		ExecutorLogger.Printf("Category with ID: %v can not be added to instruction with ID: %v", categoryID, instructionID)
		return ErrCategoryCanNotBeAddedToInstruction
	}

	return nil
}
