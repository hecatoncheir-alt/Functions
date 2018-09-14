package function

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestPriceCanBeAddedToProduct(t *testing.T) {
	PriceTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: MockStorage{DatabaseGateway: ""}}

	err := executor.AddPriceToProduct(ProductTestID, PriceTestID)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) AddEntityToOtherEntity(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestPriceCanNotBeAddedToProduct(t *testing.T) {
	PriceTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: AddPriceToProductErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddPriceToProduct(ProductTestID, PriceTestID)
	if err != ErrPriceCanNotBeAddedToProduct {
		t.Fatalf(err.Error())
	}
}

type AddPriceToProductErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddPriceToProductErrorMockStorage) AddEntityToOtherEntity(subject, predicate, object string) error {
	var status error

	if predicate == "has_price" {
		status = errors.New("")
	}

	return status
}

// --------------------------------------------------------------------------------------------------------
func TestProductCanNotBeAddedToPrice(t *testing.T) {

	PriceTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: AddProductToPriceErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddPriceToProduct(ProductTestID, PriceTestID)
	if err != ErrProductCanNotBeAddedToPrice {
		t.Fatalf(err.Error())
	}
}

type AddProductToPriceErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddProductToPriceErrorMockStorage) AddEntityToOtherEntity(subject, predicate, object string) error {
	var status error

	if predicate == "belongs_to_product" {
		status = errors.New("")
	}

	return status
}
