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
}

type Functions interface {
	ReadPriceByID(string, string) storage.Price
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrPriceCanNotBeCreated means that the price can't be added to database
	ErrPriceCanNotBeCreated = errors.New("price can't be created")
)

//// CreatePrice make price and save it to storage
func (executor *Executor) CreatePrice(price storage.Price, language string) (storage.Price, error) {
	price.IsActive = true

	encodedProduct, err := json.Marshal(price)
	if err != nil {
		ExecutorLogger.Printf(ErrPriceCanNotBeCreated.Error())
		return price, ErrPriceCanNotBeCreated
	}

	uidOfCreatedPrice, err := executor.Store.CreateJSON(encodedProduct)
	if err != nil {
		return price, ErrPriceCanNotBeCreated
	}

	createdPrice := executor.Functions.ReadPriceByID(uidOfCreatedPrice, language)

	return createdPrice, nil
}
