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

func TestExecutor_ReadProductByID(t *testing.T) {
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
		productName: string @lang @index(term, trigram) .
		productIri: string @index(term) .
		productImageLink: string @index(term) .
		productIsActive: bool @index(bool) .
		belongs_to_category: uid .
		belongs_to_company: uid .
	`

	err = setUpSchema(schema, databaseClient)

	EntityID := ""
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadProductByID(EntityID, Language)
	if err != ErrProductCanNotBeWithoutID {
		t.Error(err)
	}

	FakeID := "0x12"
	_, err = executor.ReadProductByID(FakeID, Language)
	if err != ErrProductDoesNotExist {
		t.Error(err)
	}

	entityTestName := "Test product name"
	entityTestIRI := "//"

	entityForCreate := storage.Product{
		Name:     entityTestName,
		IRI:      entityTestIRI,
		IsActive: true}

	createdEntityID, err := createEntity(entityForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdEntityID == "" {
		t.Fatalf("Created entity id is empty")
	}

	err = addOtherLanguageForCompanyName(createdEntityID, entityTestName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	entityFoundedInStorage, err := executor.ReadProductByID(createdEntityID, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if entityFoundedInStorage.ID != createdEntityID {
		t.Fatalf("ID: %v of founded product in storage is not ID: %v of created product", entityFoundedInStorage.ID, createdEntityID)
	}

	if entityFoundedInStorage.Name != entityTestName {
		t.Fatalf("Name: '%v' of founded product in storage is not value: '%v' of created product", entityFoundedInStorage.Name, entityTestName)
	}

	if entityFoundedInStorage.IRI != entityTestIRI {
		t.Fatalf("IRI: '%v' of founded product in storage is not DateTime: '%v' of created product", entityFoundedInStorage.IRI, entityTestIRI)
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

func createEntity(entityForCreate storage.Product, databaseClient *dataBaseClient.Dgraph) (string, error) {
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

func addOtherLanguageForCompanyName(entityID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	forEntityNamePredicate := fmt.Sprintf(`<%s> <productName> %s .`, entityID, "\""+name+"\""+"@"+language)

	mutation := &dataBaseAPI.Mutation{
		SetNquads: []byte(forEntityNamePredicate),
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	_, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return err
	}

	return nil
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
