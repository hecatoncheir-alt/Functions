package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestInstructionCanBeReadByID(t *testing.T) {
	IDOfTestedInstruction := "0x12"

	executor := Executor{Store: MockStore{}}

	instructionFromStore, err := executor.ReadInstructionByID(IDOfTestedInstruction, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if instructionFromStore.ID != "0x12" {
		t.Fatalf("Expected id of instruction: '0x12', actual: %v", instructionFromStore.ID)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "instructions":[
			  {
				"uid" : "0x12",
				"instructionLanguage": "ru",
				"instructionIsActive": true,
				"has_page": [],
				"has_city": [],
				"has_company": [],
				"has_category": []
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestInstructionCanBeReadByIDWithError(t *testing.T) {
	IDOfTestedInstruction := "0x12"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadInstructionByID(IDOfTestedInstruction, "ru")
	if err != ErrInstructionByIDCanNotBeFound {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStore struct {
	storage.Store
}

func (store ErrorMockStore) Query(request string) (response []byte, err error) {
	return []byte(""), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestInstructionCanBeReadByIDAndItCanBeEmpty(t *testing.T) {
	IDOfTestedInstruction := "0x12"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadInstructionByID(IDOfTestedInstruction, "ru")
	if err != ErrInstructionDoesNotExist {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "instructions":[]
		}
	`

	return []byte(resp), nil
}
