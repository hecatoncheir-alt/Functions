package function

import (
	"fmt"
	"errors"
	"github.com/hecatoncheir/Configuration"
	"github.com/hecatoncheir/Storage"
	"encoding/json"
)

type Request struct {
	Company  storage.Company
	Language string
}

// Handle a serverless request
func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)
		fmt.Println(warning)
	}


	return fmt.Sprintf("Hello, Go. You said: %s", string(req))
}

var (
	// ErrCompaniesByNameNotFound means than the companies does not exist in database
	ErrCompaniesByNameNotFound = errors.New("companies by name not found")

	// ErrCompanyCanNotBeCreated means that the company can't be added to database
	ErrCompanyCanNotBeCreated = errors.New("company can't be created")

	// ErrCompanyAlreadyExist means that the company is in the database already
	ErrCompanyAlreadyExist = errors.New("company already exist")
)

func CreateCompany(company storage.Company, language string) (storage.Company, error) {

	config := configuration.New()
	store := storage.New(config.Production.Database.Host, config.Production.Database.Port)
	store.Client.NewTxn()

}
