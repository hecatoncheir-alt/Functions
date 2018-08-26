package function

import (
	"errors"
	"log"
	"os"
)

type Storage interface {
	SetNQuads(string, string, string) error
}

type Executor struct {
	Store Storage
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCategoryCanNotBeAddedToProduct means that the category can't be added to product
	ErrCategoryCanNotBeAddedToProduct = errors.New("category can not be added to product")

	// ErrProductCanNotBeAddedToCategory means that the product can't be added to category
	ErrProductCanNotBeAddedToCategory = errors.New("product can not be added to category")
)

// AddCategoryToProduct method for set quad of predicate about product and category
func (executor *Executor) AddCategoryToProduct(productID, categoryID string) error {
	err := executor.Store.SetNQuads(categoryID, "has_product", productID)
	if err != nil {
		ExecutorLogger.Printf("Product with ID: %v can not be added to category with ID: %v", productID, categoryID)
		return ErrProductCanNotBeAddedToCategory
	}

	err = executor.Store.SetNQuads(productID, "belongs_to_category", categoryID)
	if err != nil {
		ExecutorLogger.Printf("Category with ID: %v can not be added to product with ID: %v", categoryID, productID)
		return ErrCategoryCanNotBeAddedToProduct
	}

	return nil
}
