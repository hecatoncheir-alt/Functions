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
	// ErrProductsByNameNotFound means than the products does not exist in database
	ErrProductsByNameNotFound = errors.New("products by name not found")

	// ErrProductsByNameCanNotBeFound means that the products can't be found in database
	ErrProductsByNameCanNotBeFound = errors.New("products by name can not be found")
)

// ReadProductsByName is a method for get all nodes by product name
func (executor *Executor) ReadProductsByName(productName, language string) ([]storage.Product, error) {

	variables := struct {
		ProductName string
		Language    string
	}{
		ProductName: productName,
		Language:    language}

	queryTemplate, err := template.New("ReadProductsByName").Parse(`{
				products(func: regexp(productName@{{.Language}}, /{{.ProductName}}/)) 
				@filter(eq(productIsActive, true) AND has(productName)) {
					uid
					productName: productName@{{.Language}}
					productIri
					previewImageLink
					productIsActive
					belongs_to_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
								belong_to_company @filter(eq(companyIsActive, true)){
									uid
									companyName: companyName@{{.Language}}
									companyIsActive
								}
							}
						}
					}
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIri
						companyIsActive
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
							belong_to_company @filter(eq(companyIsActive, true)){
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_price @filter(eq(priceIsActive, true)) {
						uid
						priceValue
						priceDateTime
						priceIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIri
							companyIsActive
						}
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
								priceIsActive
								belongs_to_company @filter(eq(companyIsActive, true)) {
									uid
									companyName: companyName@{{.Language}}
									companyIri
									companyIsActive
								}
							}
						}
						belongs_to_city @filter(eq(cityIsActive, true)) {
							uid
							cityName: cityName@{{.Language}}
							cityIsActive
						}
					}
				}
			}`)

	if err != nil {
		logger.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		logger.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		logger.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	type productsInStorage struct {
		AllProductsFoundedByName []storage.Product `json:"products"`
	}

	var foundedProducts productsInStorage
	err = json.Unmarshal(response, &foundedProducts)
	if err != nil {
		logger.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	if len(foundedProducts.AllProductsFoundedByName) == 0 {
		return nil, ErrProductsByNameNotFound
	}

	return foundedProducts.AllProductsFoundedByName, nil
}
