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

func TestExecutor_ReadCityByID(t *testing.T) {
	t.Skip("Database must be started")

	DatabaseGateway := "localhost:9080"
	databaseClient, err := connectToDatabase(DatabaseGateway)
	if err != nil {
		t.Fatalf(err.Error())
	}

	schema := `
		cityName: string @lang @index(term) .
		cityIsActive: bool @index(bool) .
	`

	err = setUpCitySchema(schema, databaseClient)

	CityID := ""
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadCityByID(CityID, Language)
	if err != ErrCityCanNotBeWithoutID {
		t.Error(err)
	}

	FakeCityID := "0x12"
	_, err = executor.ReadCityByID(FakeCityID, Language)
	if err != ErrCityDoesNotExist {
		t.Error(err)
	}

	CityName := "Test city name"

	cityForCreate := storage.City{
		Name:     CityName,
		IsActive: true}

	createdCityID, err := createCity(cityForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCityID == "" {
		t.Fatalf("Created city id is empty")
	}

	err = addOtherLanguageForCityName(createdCityID, CityName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cityFoundedInStorage, err := executor.ReadCityByID(createdCityID, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if cityFoundedInStorage.ID != createdCityID {
		t.Fatalf("ID: %v of founded city in storage is not ID: %v of created city", cityFoundedInStorage.ID, createdCityID)
	}

	err = deleteCityByID(createdCityID, databaseClient)
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

func setUpCitySchema(schema string, databaseClient *dataBaseClient.Dgraph) error {
	operation := &dataBaseAPI.Operation{Schema: schema}

	err := databaseClient.Alter(context.Background(), operation)
	if err != nil {
		return err
	}

	return nil
}

func createCity(cityForCreate storage.City, databaseClient *dataBaseClient.Dgraph) (string, error) {
	encodedCity, err := json.Marshal(cityForCreate)
	if err != nil {
		return "", err
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedCity,
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return "", nil
	}

	uid := assigned.Uids["blank-0"]

	return uid, nil
}

func addOtherLanguageForCityName(cityID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	forCityNamePredicate := fmt.Sprintf(`<%s> <cityName> %s .`, cityID, "\""+name+"\""+"@"+language)

	mutation := &dataBaseAPI.Mutation{
		SetNquads: []byte(forCityNamePredicate),
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	_, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return err
	}

	return nil
}

func deleteCityByID(cityID string, databaseClient *dataBaseClient.Dgraph) error {
	deleteCityData, err := json.Marshal(map[string]string{"uid": cityID})
	if err != nil {
		return err
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteCityData,
		CommitNow:  true}

	transaction := databaseClient.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}
