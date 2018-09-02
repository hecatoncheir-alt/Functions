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
	// ErrEntityCanNotBeWithoutID means that entity can't be found in storage for make some operation
	ErrEntityCanNotBeWithoutID = errors.New("entity can not be without id")

	// ErrEntityByIDCanNotBeDeleted means that the entity can't be deleted from database
	ErrEntityByIDCanNotBeDeleted = errors.New("entity by id can not be deleted")
)

// DeleteEntityByID is a method for delete city by ID
func (executor *Executor) DeleteEntityByID(entityID string) error {

	if entityID == "" {
		ExecutorLogger.Println(ErrEntityCanNotBeWithoutID)
		return ErrEntityCanNotBeWithoutID
	}

	deleteEntityData, err := json.Marshal(map[string]string{"uid": entityID})
	if err != nil {
		return err
	}

	err = executor.Store.DeleteJSON(deleteEntityData)
	if err != nil {
		ExecutorLogger.Println(err)
		return ErrEntityByIDCanNotBeDeleted
	}

	return nil
}
