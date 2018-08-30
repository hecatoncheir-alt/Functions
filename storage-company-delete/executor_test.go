package function

import (
	"github.com/hecatoncheir/Storage"
	"testing"
)

// ---------------------------------------------------------------------------------------------------------------------
func TestCompanyCanBeDeletedByID(t *testing.T) {
	IDOfTestedCompany := "0x12"

	executor := Executor{Store: MockStore{}}

	err := executor.DeleteCompanyByID(IDOfTestedCompany)
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
