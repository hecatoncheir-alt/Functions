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
	// ErrCompanyCanNotBeAddedToInstruction means that the company can't be added to instruction
	ErrCompanyCanNotBeAddedToInstruction = errors.New("company can not be added to instruction")
)

// AddCompanyToInstruction method for set quad of predicate about product and category
func (executor *Executor) AddCompanyToInstruction(instructionID, companyID string) error {
	err := executor.Store.AddEntityToOtherEntity(instructionID, "has_company", companyID)
	if err != nil {
		ExecutorLogger.Printf("Company with ID: %v can not be added to instruction with ID: %v", companyID, instructionID)
		return ErrCompanyCanNotBeAddedToInstruction
	}

	return nil
}
