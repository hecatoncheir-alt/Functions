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
	// ErrCompanyCanNotBeWithoutID means that company can't be found in storage for make some operation
	ErrCompanyCanNotBeWithoutID = errors.New("company can not be without id")

	// ErrCompanyByIDCanNotBeDeleted means that the company can't be deleted from database
	ErrCompanyByIDCanNotBeDeleted = errors.New("company by id can not be deleted")
)

// DeleteCompanyByID is a method for delete city by ID
func (executor *Executor) DeleteCompanyByID(companyID string) error {

	if companyID == "" {
		ExecutorLogger.Println(ErrCompanyCanNotBeWithoutID)
		return ErrCompanyCanNotBeWithoutID
	}

	deleteCompanyData, err := json.Marshal(map[string]string{"uid": companyID})
	if err != nil {
		return err
	}

	err = executor.Store.DeleteJSON(deleteCompanyData)
	if err != nil {
		ExecutorLogger.Println(err)
		return ErrCompanyByIDCanNotBeDeleted
	}

	return nil
}
