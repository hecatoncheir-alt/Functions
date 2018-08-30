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
	// ErrCityCanNotBeAddedToInstruction means that the City can't be added to instruction
	ErrCityCanNotBeAddedToInstruction = errors.New("city can not be added to instruction")
)

// AddCityToInstruction method for set quad of predicate about City and Instruction
func (executor *Executor) AddCityToInstruction(instructionID, cityID string) error {
	err := executor.Store.AddEntityToOtherEntity(instructionID, "has_city", cityID)
	if err != nil {
		ExecutorLogger.Printf("City with ID: %v can not be added to instruction with ID: %v", cityID, instructionID)
		return ErrCityCanNotBeAddedToInstruction
	}

	return nil
}
