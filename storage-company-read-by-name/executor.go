package function

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
	"text/template"
)

type Storage interface {
	Query(string) ([]byte, error)
}

type Executor struct {
	Store Storage
}

var logger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrCompaniesByNameCanNotBeFound means that the companies can't be found in database
	ErrCompaniesByNameCanNotBeFound = errors.New("companies by name can not be found")

	// ErrCompaniesByNameNotFound means than the companies does not exist in database
	ErrCompaniesByNameNotFound = errors.New("companies by name not found")
)

// ReadCompaniesByName is a method for get all nodes by categories name
func (executor *Executor) ReadCompaniesByName(companyName, language, databaseGateway string) ([]storage.Company, error) {
	variables := struct {
		CompanyName string
		Language    string
	}{
		CompanyName: companyName,
		Language:    language}

	queryTemplate, err := template.New("ReadCompaniesByName").Parse(`{
				companies(func: eq(companyName@{{.Language}}, "{{.CompanyName}}")) @filter(eq(companyIsActive, true)) {
					uid
					companyName: companyName@{{.Language}}
					companyIri
					companyIsActive
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIsActive
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
							}
						}
						has_product @filter(eq(productIsActive, true)) { #TODO: belongs_to_company mast be an companyID
							uid
							productName: productName@{{.Language}}
							productIri
							previewImageLink
							productIsActive
							belongs_to_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
							}
							belongs_to_company @filter(eq(companyIsActive, true)) {
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceCity
								priceIsActive
								belongs_to_product @filter(eq(productIsActive, true)) {
									uid
									productName: productName@{{.Language}}
									productIri
									previewImageLink
									productIsActive
									has_price @filter(eq(priceIsActive, true)) {
										uid
										priceValue
										priceDateTime
										priceCity
										priceIsActive
									}
								}
								belongs_to_city @filter(eq(cityIsActive, true)) {
									uid
									cityName: cityName@{{.Language}}
									cityIsActive
								}
							}
						}
					}
				}
			}`)

	if err != nil {
		logger.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		logger.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	type companiesInStorage struct {
		AllCompaniesFoundedByName []storage.Company `json:"companies"`
	}

	var foundedCompanies companiesInStorage
	err = json.Unmarshal(response, &foundedCompanies)
	if err != nil {
		logger.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	if len(foundedCompanies.AllCompaniesFoundedByName) == 0 {
		logger.Println(ErrCompaniesByNameNotFound)
		return nil, ErrCompaniesByNameNotFound
	}

	return foundedCompanies.AllCompaniesFoundedByName, nil
}
