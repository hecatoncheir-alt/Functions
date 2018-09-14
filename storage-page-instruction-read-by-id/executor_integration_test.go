package function

import (
	"context"
	"encoding/json"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"os"
	"testing"
)

func TestExecutor_ReadInstructionByID(t *testing.T) {
	t.Skip("Database must be started")

	DatabaseGateway := os.Getenv("DatabaseGateway")
	if DatabaseGateway == "" {
		DatabaseGateway = "localhost:9080"
	}

	databaseClient, err := connectToDatabase(DatabaseGateway)
	if err != nil {
		t.Fatalf(err.Error())
	}

	schema := `
		path: string @index(term) .
		pageInPaginationSelector: string @index(term) .
		previewImageOfSelector: string @index(term) .
		pageParamPath: string @index(term) .
		pageCityPath: string @index(term) .
		itemSelector: string @index(term) .
		nameOfItemSelector: string @index(term) .
		priceOfItemSelector: string @index(term) .
		cityInCookieKey: string @index(term) .
		cityIdForCookie: string @index(term) .
	`

	err = setUpSchema(schema, databaseClient)

	EmptyID := ""

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadPageInstructionByID(EmptyID)
	if err != ErrPageInstructionCanNotBeWithoutID {
		t.Error(err)
	}

	FakeID := "0x12"
	_, err = executor.ReadPageInstructionByID(FakeID)
	if err != ErrPageInstructionDoesNotExist {
		t.Error(err)
	}

	entityForCreate := storage.PageInstruction{Path: "//"}

	createdEntityID, err := createEntity(entityForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdEntityID == "" {
		t.Fatalf("Created entity id is empty")
	}

	entityFoundedInStorage, err := executor.ReadPageInstructionByID(createdEntityID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if entityFoundedInStorage.ID != createdEntityID {
		t.Fatalf("ID: %v of founded entity in storage is not ID: %v of created entity", entityFoundedInStorage.ID, createdEntityID)
	}

	err = deleteEntityByID(createdEntityID, databaseClient)
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

func setUpSchema(schema string, databaseClient *dataBaseClient.Dgraph) error {
	operation := &dataBaseAPI.Operation{Schema: schema}

	err := databaseClient.Alter(context.Background(), operation)
	if err != nil {
		return err
	}

	return nil
}

func createEntity(entityForCreate storage.PageInstruction, databaseClient *dataBaseClient.Dgraph) (string, error) {
	encodedEntity, err := json.Marshal(entityForCreate)
	if err != nil {
		return "", err
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedEntity,
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return "", nil
	}

	uid := assigned.Uids["blank-0"]

	return uid, nil
}

func deleteEntityByID(entityID string, databaseClient *dataBaseClient.Dgraph) error {
	deleteEntityData, err := json.Marshal(map[string]string{"uid": entityID})
	if err != nil {
		return err
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteEntityData,
		CommitNow:  true}

	transaction := databaseClient.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}
