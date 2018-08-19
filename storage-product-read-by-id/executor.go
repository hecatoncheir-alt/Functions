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

var ExecutorLogger = log.New(os.Stdout, "Executor: ", log.Lshortfile)

var (
	// ErrProductCanNotBeWithoutID means that product can't be found in storage for make some operation
	ErrProductCanNotBeWithoutID = errors.New("product can not be without id")

	// ErrProductByIDCanNotBeFound means that the product can't be found in database
	ErrProductByIDCanNotBeFound = errors.New("product by id can not be found")

	// ErrProductDoesNotExist means than the product does not exist in database
	ErrProductDoesNotExist = errors.New("product does not exist")
)

// ReadProductByID is a method for get all nodes of categories by ID
func (executor *Executor) ReadProductByID(productID, language string) (storage.Product, error) {
	product := storage.Product{}

	if productID == "" {
		ExecutorLogger.Println("Product can't be without ID")
		return product, ErrProductCanNotBeWithoutID
	}

	variables := struct {
		ProductID string
		Language  string
	}{
		ProductID: productID,
		Language:  language}

	queryTemplate, err := template.New("ReadProductByID").Parse(`{
				products(func: uid("{{.ProductID}}")) @filter(has(productName)) {
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
					has_price @filter(eq(priceIsActive, true)) (orderasc: priceDateTime) {
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
			}`)

	product.ID = productID

	if err != nil {
		log.Println(err)
		return product, ErrProductByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return product, ErrProductByIDCanNotBeFound
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		log.Println(err)
		return product, ErrProductByIDCanNotBeFound
	}

	type productsInStore struct {
		Products []storage.Product `json:"products"`
	}

	var foundedProducts productsInStore

	err = json.Unmarshal(response, &foundedProducts)
	if err != nil {
		log.Println(err)
		return product, ErrProductByIDCanNotBeFound
	}

	if len(foundedProducts.Products) == 0 {
		return product, ErrProductDoesNotExist
	}

	return foundedProducts.Products[0], nil
}
