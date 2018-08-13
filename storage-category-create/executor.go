package function

import (
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
)

type Storage interface {
	Mutate([]byte) (string, error)
	SetNQuads(string, string, string) error
}

type Functions interface {
	ReadCategoriesByName(string, string) []storage.Category
	ReadCategoryByID(string, string) storage.Category
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCategoryCanNotBeCreated means that the category can't be added to database
	ErrCategoryCanNotBeCreated = errors.New("category can't be created")

	// ErrCategoryAlreadyExist means that the category is in the database already
	ErrCategoryAlreadyExist = errors.New("category already exist")
)

//// CreateCategory make category and save it to storage
func (executor *Executor) CreateCategory(category storage.Category, language string) (storage.Category, error) {
	existsCategories := executor.Functions.ReadCategoriesByName(category.Name, language)
	if len(existsCategories) > 0 {
		ExecutorLogger.Printf("Category with name: %v exist: %v", category.Name, existsCategories[0])
		return existsCategories[0], ErrCategoryAlreadyExist
	}

	category.IsActive = true

	encodedCategory, err := json.Marshal(category)
	if err != nil {
		return category, ErrCategoryCanNotBeCreated
	}

	uidOfCreatedCategory, err := executor.Store.Mutate(encodedCategory)
	if err != nil {
		return category, ErrCategoryCanNotBeCreated
	}

	err = executor.Store.SetNQuads(uidOfCreatedCategory, "categoryName", "\""+category.Name+"\""+"@"+language)
	if err != nil {
		return category, ErrCategoryCanNotBeCreated
	}

	createdCategory := executor.Functions.ReadCategoryByID(uidOfCreatedCategory, language)

	return createdCategory, nil
}
