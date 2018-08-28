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

func TestExecutor_DeleteCityByID(t *testing.T) {
	//t.Skip("Database must be started")

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

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	err = executor.DeleteCityByID(CityID)
	if err != ErrCityCanNotBeWithoutID {
		t.Fatalf(err.Error())
	}

	FakeCityID := "0x12"

	err = executor.DeleteCityByID(FakeCityID)
	if err != nil {
		t.Fatalf(err.Error())
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

	Language := "ru"
	err = addOtherLanguageForCityName(createdCityID, CityName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cityFromStore, err := readCityByID(createdCityID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if cityFromStore.ID == "" {
		t.Fatalf("Created city not founded by id")
	}

	if cityFromStore.ID != createdCityID {
		t.Fatalf("Founded city id: %v not created city id: %v", cityFromStore.ID, createdCityID)
	}

	err = executor.DeleteCityByID(createdCityID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cityFromStore, err = readCityByID(createdCityID, Language, databaseClient)
	if err.Error() != "city by id not found" {
		t.Fatalf(err.Error())
	}

	err = deleteCityByID(createdCityID, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	cityFromStore, err = readCityByID(createdCityID, Language, databaseClient)
	if err.Error() != "city by id not found" {
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
	encodedCategory, err := json.Marshal(cityForCreate)
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

func readCityByID(cityID, language string, databaseClient *dataBaseClient.Dgraph) (storage.City, error) {

	var (
		ErrCityByIDCanNotBeFound = errors.New("city by id can not be found")

		ErrCityDoesNotExist = errors.New("city by id not found")
	)

	variables := struct {
		CityID   string
		Language string
	}{
		CityID:   cityID,
		Language: language}

	queryTemplate, err := template.New("ReadCityByID").Parse(`{
				cities(func: uid("{{.CityID}}")) @filter(has(cityName)) {
					uid
					cityName: cityName@{{.Language}}
					cityIsActive
				}
			}`)

	city := storage.City{ID: cityID}

	if err != nil {
		log.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	type citiesInStore struct {
		Cities []storage.City `json:"cities"`
	}

	var foundedCities citiesInStore

	err = json.Unmarshal(response.GetJson(), &foundedCities)
	if err != nil {
		log.Println(err)
		return city, ErrCityByIDCanNotBeFound
	}

	if len(foundedCities.Cities) == 0 {
		return city, ErrCityDoesNotExist
	}

	return foundedCities.Cities[0], nil
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
