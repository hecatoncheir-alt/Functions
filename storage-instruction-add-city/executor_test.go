package function

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCityCanBeAddedToInstruction(t *testing.T) {
	CityTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: MockStorage{DatabaseGateway: ""}}

	err := executor.AddCityToInstruction(InstructionTestID, CityTestID)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) SetNQuads(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestCityCanNotBeAddedToInstruction(t *testing.T) {
	CityTestID := "0x12"
	InstructionTestID := "0x13"

	executor := Executor{
		Store: AddCityToInstructionErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddCityToInstruction(InstructionTestID, CityTestID)
	if err != ErrCityCanNotBeAddedToInstruction {
		t.Fatalf(err.Error())
	}
}

type AddCityToInstructionErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddCityToInstructionErrorMockStorage) SetNQuads(subject, predicate, object string) error {
	var status error

	if predicate == "has_city" {
		status = errors.New("")
	}

	return status
}
