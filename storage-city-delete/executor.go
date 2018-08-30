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
	// ErrCityCanNotBeWithoutID means that city can't be found in storage for make some operation
	ErrCityCanNotBeWithoutID = errors.New("city can not be without id")

	// ErrCityByIDCanNotBeDeleted means that the city can't be deleted from database
	ErrCityByIDCanNotBeDeleted = errors.New("city by id can not be deleted")
)

// DeleteCityByID is a method for delete city by ID
func (executor *Executor) DeleteCityByID(cityID string) error {

	if cityID == "" {
		ExecutorLogger.Println(ErrCityCanNotBeWithoutID)
		return ErrCityCanNotBeWithoutID
	}

	deleteCityData, err := json.Marshal(map[string]string{"uid": cityID})
	if err != nil {
		return err
	}

	err = executor.Store.DeleteJSON(deleteCityData)
	if err != nil {
		ExecutorLogger.Println(err)
		return ErrCityByIDCanNotBeDeleted
	}

	return nil
}
