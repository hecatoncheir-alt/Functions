package function

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCategoryCanBeAddedToInstruction(t *testing.T) {
	CategoryTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: MockStorage{DatabaseGateway: ""}}

	err := executor.AddCategoryToInstruction(InstructionTestID, CategoryTestID)
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

func TestCategoryCanNotBeAddedToInstruction(t *testing.T) {
	CategoryTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: AddCategoryToInstructionErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddCategoryToInstruction(InstructionTestID, CategoryTestID)
	if err != ErrCategoryCanNotBeAddedToInstruction {
		t.Fatalf(err.Error())
	}
}

type AddCategoryToInstructionErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddCategoryToInstructionErrorMockStorage) AddEntityToOtherEntity(entityID, field, addedEntityID string) error {
	var status error

	if field == "has_category" {
		status = errors.New("")
	}

	return status
}
