package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestCategoryCanBeReadByID(t *testing.T) {
	IDOfTestedCategory := "0x12"

	executor := Executor{Store: MockStore{}}

	err := executor.DeleteCategoryByID(IDOfTestedCategory)
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
