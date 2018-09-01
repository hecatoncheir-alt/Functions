package function

import (
	"errors"
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestProductCanBeCreated(t *testing.T) {
	productForCreate := storage.Product{ID: "0x12", Name: "Test product", IsActive: true}

	executor := Executor{
		Functions: EmptyProductsFAASFunctions{FunctionsGateway: ""},
		Store:     MockStorage{DatabaseGateway: ""}}

	createdProduct, err := executor.CreateProduct(productForCreate, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdProduct.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", productForCreate.ID, createdProduct.ID)
	}

	if createdProduct.Name != "Test product" {
		t.Errorf("Expect: %v, but got: %v", productForCreate.Name, createdProduct.Name)
	}
}

/// Mock FAAS functions
type EmptyProductsFAASFunctions struct {
	FunctionsGateway string
}

func (functions EmptyProductsFAASFunctions) ReadProductsByName(productName, language string) []storage.Product {
	return []storage.Product{}
}

func (functions EmptyProductsFAASFunctions) ReadProductByID(productID, language string) storage.Product {
	return storage.Product{ID: productID, Name: "Test product"}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "0x12", nil
}

func (store MockStorage) AddLanguage(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestProductCanNotBeCreated(t *testing.T) {
	productForCreate := storage.Product{ID: "0x12", Name: "Test product", IsActive: true}

	executor := Executor{
		Functions: EmptyProductsFAASFunctions{FunctionsGateway: ""},
		Store:     ErrorMockStorage{DatabaseGateway: ""}}

	_, err := executor.CreateProduct(productForCreate, "ru")
	if err != ErrProductCanNotBeCreated {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStorage struct {
	DatabaseGateway string
}

func (store ErrorMockStorage) CreateJSON(setJson []byte) (uid string, err error) {
	return "", errors.New("")
}

func (store ErrorMockStorage) AddLanguage(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------
func TestCreatingProductCanBeExists(t *testing.T) {

	productForCreate := storage.Product{ID: "0x12", Name: "Test product", IsActive: true}

	executor := Executor{Functions: NotEmptyProductsFAASFunctions{}}

	existFirstTestCategory, err := executor.CreateProduct(productForCreate, "ru")
	if err != ErrProductAlreadyExist {
		t.Fatalf(err.Error())
	}

	if existFirstTestCategory.Name != "First test product name" {
		t.Errorf("Expect: %v, but got: %v", "First test product name", existFirstTestCategory.Name)
	}

}

type NotEmptyProductsFAASFunctions struct{}

func (functions NotEmptyProductsFAASFunctions) ReadProductsByName(productName, language string) []storage.Product {
	return []storage.Product{
		{ID: "0x12", Name: "First test product name", IsActive: true},
		{ID: "0x13", Name: "Second test product name", IsActive: true}}
}

func (functions NotEmptyProductsFAASFunctions) ReadProductByID(categoryID, language string) storage.Product {
	return storage.Product{ID: "0x13", Name: "Second test product name"}
}

// ---------------------------------------------------------------------------------------------------------------
