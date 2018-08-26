package function

import (
	"context"
	"encoding/json"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"testing"
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

	categoryForCreate := storage.Category{
		Name:     "Test category name",
		IsActive: true}

	createdCategoryID, err := createCategory(categoryForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCategoryID == "" {
		t.Fatalf("Created category id is empty")
	}

	err = executor.DeleteCategoryByID(FakeCategoryID)
	if err != nil {
		t.Fatalf(err.Error())
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
