package function

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"log"
	"testing"
	"text/template"
)

func TestExecutor_DeleteCategoryByID(t *testing.T) {
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

	CategoryID := ""

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	err = executor.DeleteCategoryByID(CategoryID)
	if err != ErrCategoryCanNotBeWithoutID {
		t.Fatalf(err.Error())
	}

	FakeCategoryID := "0x12"

	err = executor.DeleteCategoryByID(FakeCategoryID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	CategoryName := "Test category name"

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

	Language := "ru"
	err = addOtherLanguageForCategoryName(createdCategoryID, CategoryName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	categoryFromStore, err := readCategoryByID(createdCategoryID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if categoryFromStore.ID == "" {
		t.Fatalf("Created category not founded by id")
	}

	if categoryFromStore.ID != createdCategoryID {
		t.Fatalf("Founded category id: %v not created category id: %v", categoryFromStore.ID, createdCategoryID)
	}

	err = executor.DeleteCategoryByID(createdCategoryID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	categoryFromStore, err = readCategoryByID(createdCategoryID, Language, databaseClient)
	if err.Error() != "categories by id not found" {
		t.Fatalf(err.Error())
	}

	err = deleteCategoryByID(createdCategoryID, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	categoryFromStore, err = readCategoryByID(createdCategoryID, Language, databaseClient)
	if err.Error() != "categories by id not found" {
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

func readCategoryByID(categoryID, language string, databaseClient *dataBaseClient.Dgraph) (storage.Category, error) {

	var (
		ErrCategoryByIDCanNotBeFound = errors.New("category by id can not be found")

		ErrCategoryDoesNotExist = errors.New("categories by id not found")
	)

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

	category := storage.Category{ID: categoryID}

	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	type categoriesInStore struct {
		Categories []storage.Category `json:"categories"`
	}

	var foundedCategories categoriesInStore

	err = json.Unmarshal(response.GetJson(), &foundedCategories)
	if err != nil {
		log.Println(err)
		return category, ErrCategoryByIDCanNotBeFound
	}

	if len(foundedCategories.Categories) == 0 {
		return category, ErrCategoryDoesNotExist
	}

	return foundedCategories.Categories[0], nil
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
