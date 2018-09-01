package function

import (
	"errors"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCompanyCanBeAddedToProduct(t *testing.T) {
	CompanyTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: MockStorage{DatabaseGateway: ""}}

	err := executor.AddCompanyToProduct(ProductTestID, CompanyTestID)
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

func TestCompanyCanNotBeAddedToProduct(t *testing.T) {
	CompanyTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: AddCompanyToProductErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddCompanyToProduct(ProductTestID, CompanyTestID)
	if err != ErrCompanyCanNotBeAddedToProduct {
		t.Fatalf(err.Error())
	}
}

type AddCompanyToProductErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddCompanyToProductErrorMockStorage) AddEntityToOtherEntity(subject, predicate, object string) error {
	var status error

	if predicate == "belongs_to_company" {
		status = errors.New("")
	}

	return status
}

// --------------------------------------------------------------------------------------------------------
func TestProductCanNotBeAddedToCompany(t *testing.T) {

	CompanyTestID := "0x12"
	ProductTestID := "0x13"

	executor := Executor{
		Store: AddProductToCompanyErrorMockStorage{DatabaseGateway: ""}}

	err := executor.AddCompanyToProduct(ProductTestID, CompanyTestID)
	if err != ErrProductCanNotBeAddedToCompany {
		t.Fatalf(err.Error())
	}
}

type AddProductToCompanyErrorMockStorage struct {
	DatabaseGateway string
}

func (store AddProductToCompanyErrorMockStorage) AddEntityToOtherEntity(subject, predicate, object string) error {
	var status error

	if predicate == "has_product" {
		status = errors.New("")
	}

	return status
}

// ---------------------------------------------------------------------------------------------------------------
