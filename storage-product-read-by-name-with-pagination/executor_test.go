package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------

func TestProductCanBeReadByName(t *testing.T) {
	nameOfTestedProduct := "Test product"

	executor := Executor{Store: MockStore{}}

	productsFromStore, err := executor.ReadProductsByNameWithPagination(nameOfTestedProduct, "ru", 2, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(productsFromStore.Products) < 1 {
		t.Fatalf("Expected 1 product, actual: %v", len(productsFromStore.Products))
	}

	if productsFromStore.Products[0].ID != "0x12" {
		t.Fatalf("Expected id of product: '0x12', actual: %v", productsFromStore.Products[0].ID)
	}

	if productsFromStore.Products[0].Name != nameOfTestedProduct {
		t.Fatalf("Expected name of product: 'Test product', actual: %v", productsFromStore.Products[0].Name)
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) Query(request string) (response []byte, err error) {

	resp := `
		{  
		   "counters":[  
				{
					"total": 2
				}
			],

		   "products":[  
			  {  
				 "uid":"0x12",
				 "productName":"Test product",
				 "productIsActive":true,
                 "productIri": "http://",
				 "previewImageLink": "http://",
				 "belongs_to_company":[],
				 "belongs_to_category":[],
				 "has_price":[]
			  },
			  {  
				 "uid":"0x13",
				 "productName":"Other test product",
				 "productIsActive":true,
                 "productIri": "http://",
				 "previewImageLink": "http://",
				 "belongs_to_company":[],
				 "belongs_to_category":[],
				 "has_price":[]
			  }
		   ]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestCategoryCanBeReadByNameWithError(t *testing.T) {
	nameOfTestedProduct := "Test product"

	executor := Executor{Store: ErrorMockStore{}}
	_, err := executor.ReadProductsByNameWithPagination(nameOfTestedProduct, "ru", 2, 10)
	if err != ErrProductsByNameCanNotBeFound {
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
	nameOfTestedProduct := "Test product"

	executor := Executor{Store: EmptyMockStore{}}
	_, err := executor.ReadProductsByNameWithPagination(nameOfTestedProduct, "ru", 2, 10)
	if err != ErrProductsByNameNotFound {
		t.Fatalf(err.Error())
	}
}

type EmptyMockStore struct {
	storage.Store
}

func (store EmptyMockStore) Query(request string) (response []byte, err error) {

	resp := `
		{  
		   "counters":[  
				{
					"total": 0
				}
			],
		   "products":[]
		}
	`

	return []byte(resp), nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestProductCanBeReadByNameWithPagination(t *testing.T) {
	nameOfTestedProduct := "Test product"

	executor := Executor{Store: MockStoreForPagination{}}

	productsFromStore, err := executor.ReadProductsByNameWithPagination(nameOfTestedProduct, "ru", 2, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(productsFromStore.Products) < 1 {
		t.Fatalf("Expected 1 product, actual: %v", len(productsFromStore.Products))
	}

	if productsFromStore.Products[0].ID != "0x12" {
		t.Fatalf("Expected id of product: '0x12', actual: %v", productsFromStore.Products[0].ID)
	}

	if productsFromStore.Products[0].Name != nameOfTestedProduct {
		t.Fatalf("Expected name of product: 'Test product', actual: %v", productsFromStore.Products[0].Name)
	}

	if productsFromStore.CurrentPage != 2 {
		t.Fatalf("Expected current page is: 2, actual: %v", productsFromStore.CurrentPage)
	}

	if productsFromStore.SearchedName != nameOfTestedProduct {
		t.Fatalf("Expected name of searched product: 'Test product', actual: %v", productsFromStore.SearchedName)
	}

	if productsFromStore.TotalProductsForOnePage != 10 {
		t.Fatalf("Expected 10 products for one page, actual: %v", productsFromStore.TotalProductsForOnePage)
	}

	if productsFromStore.TotalProductsFound != 2 {
		t.Fatalf("Expected 2 products found, actual: %v", productsFromStore.TotalProductsFound)
	}
}

type MockStoreForPagination struct {
	storage.Store
}

func (store MockStoreForPagination) Query(request string) (response []byte, err error) {

	resp := `
		{  
		   "counters":[  
				{
					"total": 2
				}
			],

		   "products":[  
			  {  
				 "uid":"0x12",
				 "productName":"Test product",
				 "productIsActive":true,
                 "productIri": "http://",
				 "previewImageLink": "http://",
				 "belongs_to_company":[],
				 "belongs_to_category":[],
				 "has_price":[]
			  },
			  {  
				 "uid":"0x13",
				 "productName":"Other test product",
				 "productIsActive":true,
                 "productIri": "http://",
				 "previewImageLink": "http://",
				 "belongs_to_company":[],
				 "belongs_to_category":[],
				 "has_price":[]
			  }
		   ]
		}
	`

	return []byte(resp), nil
}
