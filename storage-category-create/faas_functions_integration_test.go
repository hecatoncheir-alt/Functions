package function

import (
	"context"
	"encoding/json"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"os"
	"testing"
)

func TestIntegration_FAASFunctions(t *testing.T) {
	t.Skip("Database and FAAS must be started")

	DatabaseGateway := os.Getenv("DatabaseGateway")
	if DatabaseGateway == "" {
		DatabaseGateway = "192.168.99.100:31285"
	}

	FunctionsGateway := os.Getenv("FunctionsGateway")
	if FunctionsGateway == "" {
		FunctionsGateway = "http://192.168.99.100:31112/function"
	}

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

	err = setUpSchema(schema, databaseClient)

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

	defer func() {
		err = deleteCategoryByID(createdCategoryID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	Language := "ru"

	err = addOtherLanguageForCategoryName(createdCategoryID, CategoryName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	executor := Executor{
		Store: &storage.Store{DatabaseGateway: DatabaseGateway},
		Functions: &FAASFunctions{
			DatabaseGateway:  DatabaseGateway,
			FunctionsGateway: FunctionsGateway}}

	foundedCategory := executor.Functions.ReadCategoryByID(createdCategoryID, Language)

	if foundedCategory.ID == createdCategoryID {
		t.Fatalf("Created category: %v by id not found", createdCategoryID)
	}

	if foundedCategory.Name == CategoryName {
		t.Fatalf("Created category: %v by name not found", CategoryName)
	}

	foundedCategories := executor.Functions.ReadCategoriesByName(foundedCategory.Name, Language)

	if len(foundedCategories) != 1 {
		t.Fatalf("No one category by name: %v found in database", CategoryName)
	}

	if foundedCategories[0].ID == createdCategoryID {
		t.Fatalf("Created category: %v by id not found", createdCategoryID)
	}

	if foundedCategories[0].Name == CategoryName {
		t.Fatalf("Created category: %v by name not found", CategoryName)
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

func setUpSchema(schema string, databaseClient *dataBaseClient.Dgraph) error {
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
		return "", err
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
