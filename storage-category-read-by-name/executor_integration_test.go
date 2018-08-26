package function

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"log"
	"testing"
	"text/template"
)

func TestExecutor_ReadCategoriesByName(t *testing.T) {
	t.Skip("Database must be started")

	DatabaseGateway := "localhost:9080"
	databaseClient, err := connectToDatabase(DatabaseGateway)
	if err != nil {
		t.Fatalf(err.Error())
	}

	schema := `
		categoryName: string @lang @index(term).
		categoryIsActive: bool @index(bool) .
		belongs_to_company: uid .
		has_product: uid .
	`

	err = setUpCategorySchema(schema, databaseClient)

	CategoryName := "Test category name"
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadCategoriesByName(CategoryName, Language)
	if err != ErrCategoriesByNameNotFound {
		t.Error(err)
	}

	categoryForCreate := storage.Category{
		Name:     CategoryName,
		IsActive: true}

	createdCategoryID, err := createCategory(categoryForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCategoryID == "" {
		t.Fatalf("Created category id is empty")
	}

	err = addOtherLanguageForCategoryName(createdCategoryID, CategoryName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	categoriesByNameFromDatabase, err := readCategoriesByName(CategoryName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(categoriesByNameFromDatabase) < 1 {
		t.Fatalf("No one category by name: %v found in database", CategoryName)
	}

	foundedCategories, err := executor.ReadCategoriesByName(CategoryName, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(foundedCategories) < 1 {
		t.Fatalf("No one category by name: %v found in database", CategoryName)
	}

	err = deleteCategoryByID(createdCategoryID, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func connectToDatabase(databaseGateway string) (*dataBaseClient.Dgraph, error) {
	conn, err := grpc.Dial(databaseGateway, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	baseConnection := dataBaseAPI.NewDgraphClient(conn)
	databaseClient := dataBaseClient.NewDgraphClient(baseConnection)

	return databaseClient, nil
}

func setUpCategorySchema(schema string, databaseClient *dataBaseClient.Dgraph) error {
	operation := &dataBaseAPI.Operation{Schema: schema}

	err := databaseClient.Alter(context.Background(), operation)
	if err != nil {
		return err
	}

	return nil
}

func createCategory(categoryForCreate storage.Category, databaseClient *dataBaseClient.Dgraph) (string, error) {
	encodedCategory, err := json.Marshal(categoryForCreate)
	if err != nil {
		return "", err
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedCategory,
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return "", nil
	}

	uid := assigned.Uids["blank-0"]

	return uid, nil
}

func addOtherLanguageForCategoryName(categoryID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	forCategoryNamePredicate := fmt.Sprintf(`<%s> <categoryName> %s .`, categoryID, "\""+name+"\""+"@"+language)

	mutation := &dataBaseAPI.Mutation{
		SetNquads: []byte(forCategoryNamePredicate),
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	_, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return err
	}

	return nil
}

func readCategoriesByName(categoryName, language string, databaseClient *dataBaseClient.Dgraph) ([]storage.Category, error) {

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
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	type categoriesInStorage struct {
		AllCategoriesFoundedByName []storage.Category `json:"categories"`
	}

	var foundedCategories categoriesInStorage
	err = json.Unmarshal(response.GetJson(), &foundedCategories)
	if err != nil {
		log.Println(err)
		return nil, ErrCategoriesByNameCanNotBeFound
	}

	if len(foundedCategories.AllCategoriesFoundedByName) == 0 {
		return nil, ErrCategoriesByNameNotFound
	}

	return foundedCategories.AllCategoriesFoundedByName, nil
}

func deleteCategoryByID(categoryID string, databaseClient *dataBaseClient.Dgraph) error {
	deleteCategoryData, err := json.Marshal(map[string]string{"uid": categoryID})
	if err != nil {
		return err
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteCategoryData,
		CommitNow:  true}

	transaction := databaseClient.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}
