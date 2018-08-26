package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestPageInstructionCanBeReadByID(t *testing.T) {
	IDOfTestedPageInstruction := "0x12"

	executor := Executor{Store: MockStore{}}

	pageInstructionFromStore, err := executor.ReadPageInstructionByID(IDOfTestedPageInstruction, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if pageInstructionFromStore.ID != "0x12" {
		t.Fatalf("Expected id of page instruction: '0x12', actual: %v", pageInstructionFromStore.ID)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "pageInstructions":[
			  {
				"uid" : "0x12",
				"path": "smartfony-i-svyaz/smartfony-205",
				"previewImageOfSelector": ".c-product-tile-picture__link .lazy-load-image-holder img",
				"itemSelector": ".c-product-tile",
				"nameOfItemSelector": ".c-product-tile__description .sel-product-tile-title",
				"linkOfItemSelector": ".c-product-tile__description .sel-product-tile-title",
				"priceOfItemSelector": ".c-product-tile__checkout-section .c-pdp-price__current"
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestPageInstructionCanBeReadByIDWithError(t *testing.T) {
	IDOfTestedPageInstruction := "0x12"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadPageInstructionByID(IDOfTestedPageInstruction, "ru")
	if err != ErrPageInstructionByIDCanNotBeFound {
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

func TestPageInstructionCanBeReadByIDAndItCanBeEmpty(t *testing.T) {
	IDOfTestedPageInstruction := "0x12"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadPageInstructionByID(IDOfTestedPageInstruction, "ru")
	if err != ErrPageInstructionDoesNotExist {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "pageInstructions":[]
		}
	`

	return []byte(resp), nil
}
