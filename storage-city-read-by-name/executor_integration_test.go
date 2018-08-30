package function

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"log"
	"testing"
	"text/template"
)

func TestExecutor_ReadCitiesByName(t *testing.T) {
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

	CityName := "Test city name"
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadCitiesByName(CityName, Language)
	if err != ErrCitiesByNameNotFound {
		t.Error(err)
	}

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

	citiesByNameFromDatabase, err := readCitiesByName(CityName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(citiesByNameFromDatabase) < 1 {
		t.Fatalf("No one city by name: %v found in database", CityName)
	}

	foundedCities, err := executor.ReadCitiesByName(CityName, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(foundedCities) < 1 {
		t.Fatalf("No one city by name: %v found in database", CityName)
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

func readCitiesByName(cityName, language string, databaseClient *dataBaseClient.Dgraph) ([]storage.City, error) {

	variables := struct {
		CityName string
		Language string
	}{
		CityName: cityName,
		Language: language}

	queryTemplate, err := template.New("ReadCitiesByName").Parse(`{
				cities(func: eq(cityName@{{.Language}}, "{{.CityName}}")) @filter(eq(cityIsActive, true)) {
					uid
					cityName: cityName@{{.Language}}
					cityIsActive
				}
			}`)

	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	type citiesInStorage struct {
		AllCitiesFoundedByName []storage.City `json:"cities"`
	}

	var foundedCities citiesInStorage
	err = json.Unmarshal(response.GetJson(), &foundedCities)
	if err != nil {
		log.Println(err)
		return nil, ErrCitiesByNameCanNotBeFound
	}

	if len(foundedCities.AllCitiesFoundedByName) == 0 {
		return nil, ErrCitiesByNameNotFound
	}

	return foundedCities.AllCitiesFoundedByName, nil
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
