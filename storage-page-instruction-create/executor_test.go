package function

import (
	"errors"
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestPageInstructionCanBeCreated(t *testing.T) {
	pageInstructionForCreate := storage.PageInstruction{Path: "http://"}

	executor := Executor{
		Functions: EmptyPageInstructionFAASFunctions{FunctionsGateway: ""},
		Store:     MockStorage{DatabaseGateway: ""}}

	createdPageInstruction, err := executor.CreatePageInstruction(pageInstructionForCreate, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdPageInstruction.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", pageInstructionForCreate.ID, createdPageInstruction.ID)
	}

	if createdPageInstruction.Path != "http://" {
		t.Errorf("Expect: %v, but got: %v", pageInstructionForCreate.Path, createdPageInstruction.Path)
	}
}

/// Mock FAAS functions
type EmptyPageInstructionFAASFunctions struct {
	FunctionsGateway string
}

func (functions EmptyPageInstructionFAASFunctions) ReadPageInstructionByID(pageInstructionID, language string) storage.PageInstruction {
	return storage.PageInstruction{ID: pageInstructionID, Path: "http://"}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "0x12", nil
}

// --------------------------------------------------------------------------------------------------------

func TestPageInstructionCanNotBeCreated(t *testing.T) {
	PageInstructionForCreate := storage.PageInstruction{}

	executor := Executor{
		Functions: EmptyPageInstructionFAASFunctions{FunctionsGateway: ""},
		Store:     ErrorMockStorage{DatabaseGateway: ""}}

	_, err := executor.CreatePageInstruction(PageInstructionForCreate, "ru")
	if err != ErrPageInstructionCanNotBeCreated {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStorage struct {
	DatabaseGateway string
}

func (store ErrorMockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "", errors.New("")
}
