package function

import (
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
)

type Storage interface {
	CreateJSON([]byte) (string, error)
	AddLanguage(string, string, string) error
}

type Functions interface {
	ReadProductsByName(string, string) []storage.Product
	ReadProductByID(string, string) storage.Product
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrProductCanNotBeCreated means that the product can't be added to database
	ErrProductCanNotBeCreated = errors.New("product can't be created")

	// ErrProductAlreadyExist means that the product is in the database already
	ErrProductAlreadyExist = errors.New("product already exist")
)

//// CreateProduct make product and save it to storage
func (executor *Executor) CreateProduct(product storage.Product, language string) (storage.Product, error) {
	existsProducts := executor.Functions.ReadProductsByName(product.Name, language)
	if len(existsProducts) > 0 {
		ExecutorLogger.Printf("Product with name: %v exist: %v", product.Name, existsProducts[0])
		return existsProducts[0], ErrProductAlreadyExist
	}

	product.IsActive = true

	encodedProduct, err := json.Marshal(product)
	if err != nil {
		return product, ErrProductCanNotBeCreated
	}

	uidOfCreatedProduct, err := executor.Store.CreateJSON(encodedProduct)
	if err != nil {
		return product, ErrProductCanNotBeCreated
	}

	err = executor.Store.AddLanguage(uidOfCreatedProduct, "productName", "\""+product.Name+"\""+"@"+language)
	if err != nil {
		return product, ErrProductCanNotBeCreated
	}

	createdProduct := executor.Functions.ReadProductByID(uidOfCreatedProduct, language)

	return createdProduct, nil
}
