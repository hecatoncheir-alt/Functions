package function

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"text/template"
	"time"
)

func TestExecutor_DeleteEntityByID(t *testing.T) {
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

	entityID := ""

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	err = executor.DeleteEntityByID(entityID)
	if err != ErrEntityCanNotBeWithoutID {
		t.Fatalf(err.Error())
	}

	FakeEntityID := "0x12"

	err = executor.DeleteEntityByID(FakeEntityID)
	if err != nil {
		t.Fatalf(err.Error())
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

	Language := "ru"

	entityFromStore, err := readEntityByID(createdEntityID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if entityFromStore.ID == "" {
		t.Fatalf("Created entity not founded by id")
	}

	if entityFromStore.ID != createdEntityID {
		t.Fatalf("Founded entity id: %v not created entity id: %v", entityFromStore.ID, createdEntityID)
	}

	err = executor.DeleteEntityByID(createdEntityID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	entityFromStore, err = readEntityByID(createdEntityID, Language, databaseClient)
	if err.Error() != "entity by id not found" {
		t.Fatalf(err.Error())
	}

	err = deleteEntityByID(createdEntityID, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	entityFromStore, err = readEntityByID(createdEntityID, Language, databaseClient)
	if err.Error() != "entity by id not found" {
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

func readEntityByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.Price, err error) {

	var (
		ErrEntityByIDCanNotBeFound = errors.New("entity by id can not be found")
		ErrEntityDoesNotExist      = errors.New("entity by id not found")
	)

	variables := struct {
		PriceID  string
		Language string
	}{
		PriceID:  entityID,
		Language: language}

	queryTemplate, err := template.New("ReadPriceByID").Parse(`{
				prices(func: uid("{{.PriceID}}")) @filter(has(priceValue)) {
					uid
					priceValue
					priceDateTime
					priceCity
					priceIsActive
					belongs_to_product @filter(eq(productIsActive, true)) {
						uid
						productName: productName@{{.Language}}
						productIri
						previewImageLink
						productIsActive
					}
					belongs_to_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@{{.Language}}
						cityIsActive
					}
					belongs_to_company @filter(eq(companyIsActive, true)){
						uid
						companyName: companyName@{{.Language}}
						companyIri
						companyIsActive
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
						}
					}
				}
			}`)

	entity = storage.Price{ID: entityID}

	if err != nil {
		log.Println(err)
		return entity, ErrEntityByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return entity, ErrEntityByIDCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return entity, ErrEntityByIDCanNotBeFound
	}

	type entitiesInStore struct {
		Entities []storage.Price `json:"prices"`
	}

	var foundedEntities entitiesInStore

	err = json.Unmarshal(response.GetJson(), &foundedEntities)
	if err != nil {
		log.Println(err)
		return entity, ErrEntityByIDCanNotBeFound
	}

	if len(foundedEntities.Entities) == 0 {
		return entity, ErrEntityDoesNotExist
	}

	return foundedEntities.Entities[0], nil
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
