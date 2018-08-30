package function

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestPageInstructionCanBeAddedToInstruction(t *testing.T) {
	PageInstructionTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: MockStorage{DatabaseGateway: ""}}

	err := executor.AddPageInstructionToInstruction(InstructionTestID, PageInstructionTestID)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) AddEntityToOtherEntity(entityID, field, addedEntityID string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestPageInstructionCanNotBeAddedToInstruction(t *testing.T) {
	PageInstructionTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: AddCompanyToInstructionErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddPageInstructionToInstruction(InstructionTestID, PageInstructionTestID)
	if err != ErrPageInstructionCanNotBeAddedToInstruction {
		t.Fatalf(err.Error())
	}
}

type AddCompanyToInstructionErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddCompanyToInstructionErrorMockStorage) AddEntityToOtherEntity(entityID, field, addedEntityID string) error {
	var status error

	if field == "has_page" {
		status = errors.New("")
	}

	return status
}
