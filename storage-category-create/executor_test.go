package function

import (
	"errors"
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ------------------------------------------------------------------------------------------------------
func TestCategoryCanBeCreated(t *testing.T) {
	categoryForCreate := storage.Category{ID: "0x12", Name: "Test category", IsActive: true}

	executor := Executor{
		Functions: EmptyCategoriesFAASFunctions{FunctionsGateway: ""},
		Store:     MockStorage{DatabaseGateway: ""}}

	createdCategory, err := executor.CreateCategory(categoryForCreate, "ru")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCategory.ID != "0x12" {
		t.Errorf("Expect: %v, but got: %v", categoryForCreate.ID, createdCategory.ID)
	}

	if createdCategory.Name != "Test category" {
		t.Errorf("Expect: %v, but got: %v", categoryForCreate.Name, createdCategory.Name)
	}
}

/// Mock FAAS functions
type EmptyCategoriesFAASFunctions struct {
	FunctionsGateway string
}

func (functions EmptyCategoriesFAASFunctions) ReadCategoriesByName(categoryName, language string) []storage.Category {
	return []storage.Category{}
}

func (functions EmptyCategoriesFAASFunctions) ReadCategoryByID(categoryID, language string) storage.Category {
	return storage.Category{ID: categoryID, Name: "Test category"}
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

func TestCategoryCanNotBeCreated(t *testing.T) {
	categoryForCreate := storage.Category{ID: "0x12", Name: "Test category", IsActive: true}

	executor := Executor{
		Functions: EmptyCategoriesFAASFunctions{FunctionsGateway: ""},
		Store:     ErrorMockStorage{DatabaseGateway: ""}}

	_, err := executor.CreateCategory(categoryForCreate, "ru")
	if err != ErrCategoryCanNotBeCreated {
		t.Fatalf(err.Error())
	}
}

type ErrorMockStorage struct {
	DatabaseGateway string
}

func (store ErrorMockStorage) Mutate(setJson []byte) (uid string, err error) {
	return "", errors.New("")
}

func (store ErrorMockStorage) SetNQuads(subject, predicate, object string) error {
	return nil
}

// --------------------------------------------------------------------------------------------------------
func TestCreatingCategoryCanBeExists(t *testing.T) {

	categoryForCreate := storage.Category{ID: "0x12", Name: "Test category", IsActive: true}

	executor := Executor{Functions: NotEmptyCategoriesFAASFunctions{}}

	existFirstTestCategory, err := executor.CreateCategory(categoryForCreate, "ru")
	if err != ErrCategoryAlreadyExist {
		t.Fatalf(err.Error())
	}

	if existFirstTestCategory.Name != "First test category name" {
		t.Errorf("Expect: %v, but got: %v", "First test category name", existFirstTestCategory.Name)
	}

}

type NotEmptyCategoriesFAASFunctions struct{}

func (functions NotEmptyCategoriesFAASFunctions) ReadCategoriesByName(categoryName, language string) []storage.Category {
	return []storage.Category{
		{ID: "0x12", Name: "First test category name", IsActive: true},
		{ID: "0x13", Name: "Second test category name", IsActive: true}}
}

func (functions NotEmptyCategoriesFAASFunctions) ReadCategoryByID(categoryID, language string) storage.Category {
	return storage.Category{ID: "0x13", Name: "Second test category name"}
}

// ---------------------------------------------------------------------------------------------------------------
