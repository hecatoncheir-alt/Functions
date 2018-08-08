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
	// ErrCompanyCanNotBeWithoutID means that company can't be found in storage for make some operation
	ErrCompanyCanNotBeWithoutID = errors.New("company can not be without id")

	// ErrCompanyByIDCanNotBeFound means that the company can't be found in database
	ErrCompanyByIDCanNotBeFound = errors.New("company by id can not be found")

	// ErrCompanyDoesNotExist means than the company does not exist in database
	ErrCompanyDoesNotExist = errors.New("company does not exist")
)

// ReadCompanyByID is a method for get all nodes of categories by ID
func (executor *Executor) ReadCompanyByID(companyID, language string) (storage.Company, error) {
	company := storage.Company{ID: companyID}

	if companyID == "" {
		logger.Println(ErrCompanyCanNotBeWithoutID)
		return company, ErrCompanyCanNotBeWithoutID
	}

	variables := struct {
		CompanyID string
		Language  string
	}{
		CompanyID: companyID,
		Language:  language}

	queryTemplate, err := template.New("ReadCompanyByID").Parse(`{
				companies(func: uid("{{.CompanyID}}")) @filter(has(companyName)) {
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
						has_product @filter(uid_in(belongs_to_company, {{.CompanyID}}) AND eq(productIsActive, true)) {
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
		return company, ErrCompanyByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		logger.Println(err)
		return company, err
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		logger.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	type companiesInStore struct {
		Companies []storage.Company `json:"companies"`
	}

	var foundedCompanies companiesInStore

	err = json.Unmarshal(response, &foundedCompanies)
	if err != nil {
		logger.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	if len(foundedCompanies.Companies) == 0 {
		return company, ErrCompanyDoesNotExist
	}

	return foundedCompanies.Companies[0], nil
}
