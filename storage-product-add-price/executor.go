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
	// ErrProductCanNotBeAddedToPrice means that the product can't be added to price
	ErrProductCanNotBeAddedToPrice = errors.New("product can not be added to price")

	// ErrPriceCanNotBeAddedToProduct means that the price can't be added to product
	ErrPriceCanNotBeAddedToProduct = errors.New("price can not be added to product")
)

// AddPriceToProduct method for set quad of predicate about product and category
func (executor *Executor) AddPriceToProduct(productID, priceID string) error {
	err := executor.Store.AddEntityToOtherEntity(priceID, "belongs_to_product", productID)
	if err != nil {
		ExecutorLogger.Printf("Product with ID: %v can not be added to price with ID: %v", productID, priceID)
		return ErrProductCanNotBeAddedToPrice
	}

	err = executor.Store.AddEntityToOtherEntity(productID, "has_price", priceID)
	if err != nil {
		ExecutorLogger.Printf("Price with ID: %v can not be added to product with ID: %v", priceID, productID)
		return ErrPriceCanNotBeAddedToProduct
	}

	return nil
}
