package function

import (
	"errors"
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestPriceCanBeCreated(t *testing.T) {
	priceForCreate := storage.Price{ID: "0x12", Value: 0.0, IsActive: true}

	executor := Executor{
		Functions: EmptyPriceFAASFunctions{FunctionsGateway: ""},
		Store:     MockStorage{DatabaseGateway: ""}}

	createdPrice, err := executor.CreatePrice(priceForCreate, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdPrice.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", priceForCreate.ID, createdPrice.ID)
	}

	if createdPrice.Value != 0.0 {
		t.Errorf("Expect: %v, but got: %v", priceForCreate.Value, createdPrice.Value)
	}
}

/// Mock FAAS functions
type EmptyPriceFAASFunctions struct {
	FunctionsGateway string
}

func (functions EmptyPriceFAASFunctions) ReadPriceByID(priceID, language string) storage.Price {
	return storage.Price{ID: priceID, Value: 0.0}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) Mutate(setJson []byte) (uid string, err error) {
	return "0x12", nil
}

func (store MockStorage) SetNQuads(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestPriceCanNotBeCreated(t *testing.T) {
	priceForCreate := storage.Price{ID: "0x12", Value: 0.0, IsActive: true}

	executor := Executor{
		Functions: EmptyPriceFAASFunctions{FunctionsGateway: ""},
		Store:     ErrorMockStorage{DatabaseGateway: ""}}

	_, err := executor.CreatePrice(priceForCreate, "ru")
	if err != ErrPriceCanNotBeCreated {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStorage struct {
	DatabaseGateway string
}

func (store ErrorMockStorage) Mutate(setJson []byte) (uid string, err error) {
	return "", errors.New("")
}
