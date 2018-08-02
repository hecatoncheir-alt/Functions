package function

import (
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
)

type Storage interface {
	Query(string) ([]byte, error)
}

type Functions interface {
	CompaniesReadByName(string, string, string) []storage.Company
}

type Executor struct {
	Store     Storage
	Functions Functions
}

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCompanyAlreadyExist means that the company is in the database already
	ErrCompanyAlreadyExist = errors.New("company already exist")
)

// ReadCompaniesByName is a method for get all nodes by categories name
func (executor *Executor) CreateCompany(
	company storage.Company, language, DatabaseGateway string) (storage.Company, error) {

	existsCompanies := executor.Functions.CompaniesReadByName(company.Name, language, DatabaseGateway)

	if len(existsCompanies) > 0 {
		ExecutorLogger.Printf("Company with name: %v exist: %v", company.Name, existsCompanies[0])
		return existsCompanies[0], ErrCompanyAlreadyExist
	}

	// TODO

	return storage.Company{}, nil
}
