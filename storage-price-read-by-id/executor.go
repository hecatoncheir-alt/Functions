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
	// ErrPriceCanNotBeWithoutID means that price can't be without id
	ErrPriceCanNotBeWithoutID = errors.New("price can not be without id")

	// ErrPriceByIDCanNotBeFound means that the price can't be found in database
	ErrPriceByIDCanNotBeFound = errors.New("price by id can not be found")

	// ErrPriceDoesNotExist means than the price does not exist in database
	ErrPriceDoesNotExist = errors.New("price does not exist")
)

// ReadPriceByID is a method for get all nodes of categories by ID
func (executor *Executor) ReadPriceByID(priceID, language string) (storage.Price, error) {
	price := storage.Price{}

	if priceID == "" {
		logger.Println(ErrPriceCanNotBeWithoutID)
		return price, ErrPriceCanNotBeWithoutID
	}

	variables := struct {
		PriceID  string
		Language string
	}{
		PriceID:  priceID,
		Language: language}

	queryTemplate, err := template.New("ReadPriceByID").Parse(`{
				prices(func: uid("{{.PriceID}}")) @filter(has(priceValue)) {
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
					}
					belongs_to_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@{{.Language}}
						cityIsActive
					}
					belongs_to_company @filter(eq(companyIsActive, true)){
						uid
						companyName: companyName@{{.Language}}
						companyIri
						companyIsActive
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
						}
					}
				}
			}`)

	price = storage.Price{ID: priceID}
	if err != nil {
		log.Println(err)
		return price, ErrPriceByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return price, err
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		log.Println(err)
		return price, ErrPriceByIDCanNotBeFound
	}

	type PricesInStore struct {
		Prices []storage.Price `json:"prices"`
	}

	var foundedPrices PricesInStore

	err = json.Unmarshal(response, &foundedPrices)
	if err != nil {
		log.Println(err)
		return price, ErrPriceByIDCanNotBeFound
	}

	if len(foundedPrices.Prices) == 0 {
		return price, ErrPriceDoesNotExist
	}

	return foundedPrices.Prices[0], nil
}
