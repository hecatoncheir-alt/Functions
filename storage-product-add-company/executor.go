package function

import (
	"errors"
	"log"
	"os"
)

type Storage interface {
	AddEntityToOtherEntity(string, string, string) error
}

type Executor struct {
	Store Storage
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCompanyCanNotBeAddedToProduct means that the company can't be added to product
	ErrCompanyCanNotBeAddedToProduct = errors.New("company can not be added to product")

	// ErrProductCanNotBeAddedToCompany means that the product can't be added to company
	ErrProductCanNotBeAddedToCompany = errors.New("product can not be added to company")
)

// AddCompanyToProduct method for set quad of predicate about product and category
func (executor *Executor) AddCompanyToProduct(productID, companyID string) error {
	err := executor.Store.AddEntityToOtherEntity(companyID, "has_product", productID)
	if err != nil {
		ExecutorLogger.Printf("Product with ID: %v can not be added to company with ID: %v", productID, companyID)
		return ErrProductCanNotBeAddedToCompany
	}

	err = executor.Store.AddEntityToOtherEntity(productID, "belongs_to_company", companyID)
	if err != nil {
		ExecutorLogger.Printf("Company with ID: %v can not be added to product with ID: %v", companyID, productID)
		return ErrCompanyCanNotBeAddedToProduct
	}

	return nil
}
