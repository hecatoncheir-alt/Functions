package function

import (
	"context"
	"encoding/json"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"testing"
)

func TestExecutor_ReadCategoryByID(t *testing.T) {
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
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadCategoryByID(CategoryID, Language)
	if err != ErrCategoryCanNotBeWithoutID {
		t.Error(err)
	}

	FakeCategoryID := "0x12"
	_, err = executor.ReadCategoryByID(FakeCategoryID, Language)
	if err != ErrCategoryDoesNotExist {
		t.Error(err)
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

	err = addOtherLanguageForCategoryName(createdCategoryID, CategoryName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	categoryFoundedInStorage, err := executor.ReadCategoryByID(createdCategoryID, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if categoryFoundedInStorage.ID != createdCategoryID {
		t.Fatalf("ID: %v of founded category in storage is not ID: %v of created category", categoryFoundedInStorage.ID, createdCategoryID)
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
