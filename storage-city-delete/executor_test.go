package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestCityCanBeReadByID(t *testing.T) {
	IDOfTestedCity := "0x12"

	executor := Executor{Store: MockStore{}}

	err := executor.DeleteCityByID(IDOfTestedCity)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

type MockStore struct {
	storage.Store
}

func (store MockStore) DeleteJSON(encodedJSON []byte) error {
	return nil
}
