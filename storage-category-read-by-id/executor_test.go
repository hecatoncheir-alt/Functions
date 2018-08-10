package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestCategoryCanBeReadByID(t *testing.T) {
	IDOfTestedCategory := "0x12"

	executor := Executor{Store: MockStore{}}

	categoryFromStore, err := executor.ReadCategoryByID(IDOfTestedCategory, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if categoryFromStore.ID != "0x12" {
		t.Fatalf("Expected id of category: '0x12', actual: %v", categoryFromStore.ID)
	}

	if categoryFromStore.Name != "Test category" {
		t.Fatalf("Expected name of category: 'Test category', actual: %v", categoryFromStore.Name)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "categories":[
			  {
				 "uid":"0x12",
				 "categoryName":"Test category",
				 "categoryIsActive":true,
				 "belongs_to_company": [],
				 "has_product": []
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestCategoryCanBeReadByNameWithError(t *testing.T) {
	IDOfTestedCategory := "0x12"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadCategoryByID(IDOfTestedCategory, "ru")
	if err != ErrCategoryByIDCanNotBeFound {
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

func TestCompanyCanBeReadByNameAndItCanBeEmpty(t *testing.T) {
	IDOfTestedCompany := "0x12"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadCategoryByID(IDOfTestedCompany, "ru")
	if err != ErrCategoryDoesNotExist {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{
		   "categories":[]
		}
	`

	return []byte(resp), nil
}
