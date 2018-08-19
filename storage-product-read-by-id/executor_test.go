package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestProductCanBeReadByID(t *testing.T) {
	IDOfTestedProduct := "0x12"

	executor := Executor{Store: MockStore{}}

	productFromStore, err := executor.ReadProductByID(IDOfTestedProduct, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if productFromStore.ID != "0x12" {
		t.Fatalf("Expected id of product: '0x12', actual: %v", productFromStore.ID)
	}

	if productFromStore.Name != "Test product" {
		t.Fatalf("Expected name of product: 'Test product', actual: %v", productFromStore.Name)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "products":[
			  {
				 "uid":"0x12",
				 "productName":"Test product",
				 "productIri":"http://",
				 "previewImageLink":"http://",
				 "productIsActive":true,
				 "belongs_to_company": [],
				 "belongs_to_category": [],
				 "has_price": []
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestProductCanBeReadByIDWithError(t *testing.T) {
	IDOfTestedProduct := "0x12"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadProductByID(IDOfTestedProduct, "ru")
	if err != ErrProductByIDCanNotBeFound {
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

func TestProductCanBeReadByIDAndItCanBeEmpty(t *testing.T) {
	IDOfTestedProduct := "0x12"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadProductByID(IDOfTestedProduct, "ru")
	if err != ErrProductDoesNotExist {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "products":[]
		}
	`

	return []byte(resp), nil
}
