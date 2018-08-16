package function

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCategoryCanBeAddedToProduct(t *testing.T) {
	CategoryTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: MockStorage{DatabaseGateway: ""}}

	err := executor.AddCategoryToProduct(ProductTestID, CategoryTestID)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

/// Mock Storage
type MockStorage struct {
	DatabaseGateway string
}

func (store MockStorage) SetNQuads(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------

func TestCategoryCanNotBeAddedToProduct(t *testing.T) {
	CategoryTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: AddCategoryToProductErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddCategoryToProduct(ProductTestID, CategoryTestID)
	if err != ErrCategoryCanNotBeAddedToProduct {
		t.Fatalf(err.Error())
	}
}

type AddCategoryToProductErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddCategoryToProductErrorMockStorage) SetNQuads(subject, predicate, object string) error {
	var status error

	if predicate == "belongs_to_category" {
		status = errors.New("")
	}

	return status
}

// --------------------------------------------------------------------------------------------------------
func TestProductCanNotBeAddedToCategory(t *testing.T) {

	CategoryTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: AddProductToCategoryErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddCategoryToProduct(ProductTestID, CategoryTestID)
	if err != ErrProductCanNotBeAddedToCategory {
		t.Fatalf(err.Error())
	}
}

type AddProductToCategoryErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddProductToCategoryErrorMockStorage) SetNQuads(subject, predicate, object string) error {
	var status error

	if predicate == "has_product" {
		status = errors.New("")
	}

	return status
}

// ---------------------------------------------------------------------------------------------------------------
