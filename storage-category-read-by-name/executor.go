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
	// ErrCategoriesByNameNotFound means than the categories does not exist in database
	ErrCategoriesByNameNotFound = errors.New("categories by name not found")

	// ErrCategoriesByNameCanNotBeFound means that the category can't be found in database
	ErrCategoriesByNameCanNotBeFound = errors.New("categories by name can not be found")
)

// ReadCategoriesByName is a method for get all nodes by categories name
func (executor *Executor) ReadCategoriesByName(categoryName, language string) ([]storage.Category, error) {
	variables := struct {
		CategoryName string
		Language     string
	}{
		CategoryName: categoryName,
		Language:     language}

	queryTemplate, err := template.New("ReadCategoriesByName").Parse(`{
				categories(func: eq(categoryName@{{.Language}}, "{{.CategoryName}}"))
				@filter(eq(categoryIsActive, true)) {
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

	if err != nil {
		ExecutorLogger.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	response, err := executor.Store.Query(queryBuf.String())
	if err != nil {
		ExecutorLogger.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	type categoriesInStorage struct {
		AllCategoriesFoundedByName []storage.Category `json:"categories"`
	}

	var foundedCategories categoriesInStorage
	err = json.Unmarshal(response, &foundedCategories)
	if err != nil {
		ExecutorLogger.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	if len(foundedCategories.AllCategoriesFoundedByName) == 0 {
		return nil, ErrCategoriesByNameNotFound
	}

	return foundedCategories.AllCategoriesFoundedByName, nil
}
