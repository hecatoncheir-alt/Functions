package function

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCompanyCanBeAddedToInstruction(t *testing.T) {
	CompanyTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: MockStorage{DatabaseGateway: ""}}

	err := executor.AddCompanyToInstruction(InstructionTestID, CompanyTestID)
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

func TestCompanyCanNotBeAddedToInstruction(t *testing.T) {
	CompanyTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: AddCompanyToInstructionErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddCompanyToInstruction(InstructionTestID, CompanyTestID)
	if err != ErrCompanyCanNotBeAddedToInstruction {
		t.Fatalf(err.Error())
	}
}

type AddCompanyToInstructionErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddCompanyToInstructionErrorMockStorage) AddEntityToOtherEntity(entityID, field, addedEntityID string) error {
	var status error

	if field == "has_company" {
		status = errors.New("")
	}

	return status
}
