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
	"time"
)

func TestExecutor_ReadPriceByID(t *testing.T) {
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
		pricesValue: float @index(float) .
		priceDateTime: dateTime @index(day) .
		priceIsActive: bool @index(bool) .
		belongs_to_city: uid @count .
		belongs_to_product: uid @count .
		belongs_to_company: uid @count .
	`

	err = setUpSchema(schema, databaseClient)

	EntityID := ""
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadPriceByID(EntityID, Language)
	if err != ErrPriceCanNotBeWithoutID {
		t.Error(err)
	}

	FakeID := "0x12"
	_, err = executor.ReadPriceByID(FakeID, Language)
	if err != ErrPriceDoesNotExist {
		t.Error(err)
	}

	testPriceValue := 0.1
	testPriceDateTime := time.Now().UTC()

	entityForCreate := storage.Price{
		Value:    testPriceValue,
		DateTime: testPriceDateTime,
		IsActive: true}

	createdEntityID, err := createEntity(entityForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdEntityID == "" {
		t.Fatalf("Created entity id is empty")
	}

	entityFoundedInStorage, err := executor.ReadPriceByID(createdEntityID, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if entityFoundedInStorage.ID != createdEntityID {
		t.Fatalf("ID: %v of founded price in storage is not ID: %v of created price", entityFoundedInStorage.ID, createdEntityID)
	}

	if entityFoundedInStorage.Value != testPriceValue {
		t.Fatalf("Value: %v of founded price in storage is not value: %v of created price", entityFoundedInStorage.Value, testPriceValue)
	}

	if entityFoundedInStorage.DateTime != testPriceDateTime {
		t.Fatalf("DateTime: %v of founded price in storage is not DateTime: %v of created price", entityFoundedInStorage.Value, testPriceValue)
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

func createEntity(entityForCreate storage.Price, databaseClient *dataBaseClient.Dgraph) (string, error) {
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
