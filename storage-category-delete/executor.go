package function

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Storage interface {
	DeleteJSON([]byte) error
}

type Executor struct {
	Store Storage
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCategoryCanNotBeWithoutID means that category can't be found in storage for make some operation
	ErrCategoryCanNotBeWithoutID = errors.New("category can not be without id")

	// ErrCategoryByIDCanNotBeDeleted means that the category can't be deleted from database
	ErrCategoryByIDCanNotBeDeleted = errors.New("category by id can not be deleted")
)

// DeleteCategoryByID is a method for delete category by ID
func (executor *Executor) DeleteCategoryByID(categoryID string) error {

	if categoryID == "" {
		ExecutorLogger.Println(ErrCategoryCanNotBeWithoutID)
		return ErrCategoryCanNotBeWithoutID
	}

	deleteCategoryData, err := json.Marshal(map[string]string{"uid": categoryID})
	if err != nil {
		return err
	}

	err = executor.Store.DeleteJSON(deleteCategoryData)
	if err != nil {
		ExecutorLogger.Println(err)
		return ErrCategoryByIDCanNotBeDeleted
	}

	return nil
}
