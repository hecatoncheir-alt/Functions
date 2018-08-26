package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestPriceCanBeReadByID(t *testing.T) {
	IDOfTestedPrice := "0x12"

	executor := Executor{Store: MockStore{}}

	priceFromStore, err := executor.ReadPriceByID(IDOfTestedPrice, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if priceFromStore.ID != "0x12" {
		t.Fatalf("Expected id of price: '0x12', actual: %v", priceFromStore.ID)
	}

	if priceFromStore.Value != 0.0 {
		t.Fatalf("Expected value of price: 0, actual: %v", priceFromStore.Value)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "prices":[
			  {
				"uid": "0x12",
				"priceValue": 0.0,
				"priceDateTime" : "2017-05-01T16:27:18.543653798Z",
				"priceIsActive": true,
				"belongs_to_city": [],
				"belongs_to_product": [],
				"belongs_to_company": []
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestPriceCanBeReadByIDWithError(t *testing.T) {
	IDOfTestedPrice := "0x12"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadPriceByID(IDOfTestedPrice, "ru")
	if err != ErrPriceByIDCanNotBeFound {
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

func TestPriceCanBeReadByIDAndItCanBeEmpty(t *testing.T) {
	IDOfTestedPrice := "0x12"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadPriceByID(IDOfTestedPrice, "ru")
	if err != ErrPriceDoesNotExist {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "prices":[]
		}
	`

	return []byte(resp), nil
}
