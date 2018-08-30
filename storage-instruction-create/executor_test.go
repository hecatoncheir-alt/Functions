package function

import (
	"errors"
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestInstructionCanBeCreated(t *testing.T) {
	instructionForCreate := storage.Instruction{}

	executor := Executor{
		Functions: EmptyInstructionFAASFunctions{FunctionsGateway: ""},
		Store:     MockStorage{DatabaseGateway: ""}}

	createdInstruction, err := executor.CreateInstruction(instructionForCreate, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdInstruction.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", instructionForCreate.ID, createdInstruction.ID)
	}

	if createdInstruction.IsActive != true {
		t.Errorf("Expect: %v, but got: %v", instructionForCreate.IsActive, createdInstruction.IsActive)
	}
}

/// Mock FAAS functions
type EmptyInstructionFAASFunctions struct {
	FunctionsGateway string
}

func (functions EmptyInstructionFAASFunctions) ReadInstructionByID(instructionID, language string) storage.Instruction {
	return storage.Instruction{ID: instructionID, IsActive: true}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "0x12", nil
}

// --------------------------------------------------------------------------------------------------------

func TestInstructionCanNotBeCreated(t *testing.T) {
	instructionForCreate := storage.Instruction{}

	executor := Executor{
		Functions: EmptyInstructionFAASFunctions{FunctionsGateway: ""},
		Store:     ErrorMockStorage{DatabaseGateway: ""}}

	_, err := executor.CreateInstruction(instructionForCreate, "ru")
	if err != ErrInstructionCanNotBeCreated {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStorage struct {
	DatabaseGateway string
}

func (store ErrorMockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "", errors.New("")
}
