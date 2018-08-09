package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------

func TestCompanyCanBeReadByName(t *testing.T) {
	nameOfTestedCategory := "Test category"

	executor := Executor{Store: MockStore{}}

	categoriesFromStore, err := executor.ReadCategoriesByName(nameOfTestedCategory, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(categoriesFromStore) < 1 {
		t.Fatalf("Expected 1 category, actual: %v", len(categoriesFromStore))
	}

	if categoriesFromStore[0].ID != "0x12" {
		t.Fatalf("Expected id of category: '0x12', actual: %v", categoriesFromStore[0].ID)
	}

	if categoriesFromStore[0].Name != nameOfTestedCategory {
		t.Fatalf("Expected name of category: 'Test category', actual: %v", categoriesFromStore[0].Name)
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
				 "belongs_to_company":[],
				 "has_product":[]
			  },
			  {  
				 "uid":"0x13",
				 "categoryName":"Test category",
				 "categoryIsActive":true,
				 "belongs_to_company":[],
				 "has_product":[]
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestCategoryCanBeReadByNameWithError(t *testing.T) {
	nameOfTestedCategory := "Test category"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadCategoriesByName(nameOfTestedCategory, "ru")
	if err != ErrCategoriesByNameCanNotBeFound {
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

func TestCategoryCanBeReadByNameAndItCanBeEmpty(t *testing.T) {
	nameOfTestedCategory := "Test category"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadCategoriesByName(nameOfTestedCategory, "ru")
	if err != ErrCategoriesByNameNotFound {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{  
		   "companies":[]
		}
	`

	return []byte(resp), nil
}
