package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestEntityCanBeDeletedByID(t *testing.T) {
	IDOfTestedEntity := "0x12"

	executor := Executor{Store: MockStore{}}

	err := executor.DeleteEntityByID(IDOfTestedEntity)
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
