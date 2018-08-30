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
	// ErrCategoryCanNotBeWithoutID means that category can't be found in storage for make some operation
	ErrCategoryCanNotBeWithoutID = errors.New("category can not be without id")

	// ErrCategoryByIDCanNotBeFound means that the category can't be found in database
	ErrCategoryByIDCanNotBeFound = errors.New("category by id can not be found")

	// ErrCategoryDoesNotExist means than the category does not exist in database
	ErrCategoryDoesNotExist = errors.New("category does not exist")
)

// ReadCategoryByID is a method for get all nodes of categories by ID
func (executor *Executor) ReadCategoryByID(categoryID, language string) (storage.Category, error) {
	category := storage.Category{}

	if categoryID == "" {
		ExecutorLogger.Println(ErrCategoryCanNotBeWithoutID)
		return category, ErrCategoryCanNotBeWithoutID
	}

	variables := struct {
		CategoryID string
		Language   string
	}{
		CategoryID: categoryID,
		Language:   language}

	queryTemplate, err := template.New("ReadCategoryByID").Parse(`{
				categories(func: uid("{{.CategoryID}}")) @filter(has(categoryName)) {
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
							belong_to_company @filter(eq(companyIsActive, true)) {
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_product @filter(eq(productIsActive, true)) {
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
					}
				}
			}`)

	category.ID = categoryID

	if err != nil {
		ExecutorLogger.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		ExecutorLogger.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	type categoriesInStore struct {
		Categories []storage.Category `json:"categories"`
	}

	var foundedCategories categoriesInStore

	err = json.Unmarshal(response, &foundedCategories)
	if err != nil {
		ExecutorLogger.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	if len(foundedCategories.Categories) == 0 {
		return category, ErrCategoryDoesNotExist
	}

	return foundedCategories.Categories[0], nil
}
