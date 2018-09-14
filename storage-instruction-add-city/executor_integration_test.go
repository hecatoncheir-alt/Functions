package function

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"os"
	"testing"
	"text/template"
)

func TestExecutor_AddCityToInstruction(t *testing.T) {
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
		instructionLanguage: string @index(term) .
		instructionIsActive: bool @index(bool) .
		has_company: uid @count .
		has_city: uid @count .
		has_page: uid @count .
		has_category: uid @count .
	`

	err = setUpSchema(schema, databaseClient)

	Language := "ru"

	entityForCreate := storage.Instruction{
		Language: Language,
		IsActive: true}

	createdEntityID, err := createInstruction(entityForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdEntityID == "" {
		t.Fatalf("Created entity id is empty")
	}

	defer func() {
		err = deleteEntityByID(createdEntityID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	entityFoundedInStorage, err := readInstructionByID(createdEntityID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if entityFoundedInStorage.ID != createdEntityID {
		t.Fatalf("ID: %v of founded entity in storage is not ID: %v of created entity", entityFoundedInStorage.ID, createdEntityID)
	}

	if len(entityFoundedInStorage.Cities) > 0 {
		t.Fatalf("Expect 0 cities but got: %v", entityFoundedInStorage.Cities)
	}

	cityForCreate := storage.City{
		Name:     "Test city name",
		IsActive: true}

	schema = `
		cityName: string @lang @index(term) .
		cityIsActive: bool @index(bool) .
	`

	err = setUpSchema(schema, databaseClient)

	cityID, err := createCity(cityForCreate, Language, databaseClient)
	if err != nil {
		t.Fatalf("City does not create")
	}

	defer func() {
		err = deleteEntityByID(cityID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	city, err := readCityByID(cityID, Language, databaseClient)
	if err != nil {
		t.Fatalf("City does not exist: %v", err)
	}

	if city.ID != cityID {
		t.Fatalf("ID: %v of founded city in storage is not ID: %v of created city", city.ID, cityID)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}

	err = executor.AddCityToInstruction(createdEntityID, cityID)
	if err != nil {
		t.Fatalf("Name of City does not added")
	}

	entityFoundedInStorage, err = readInstructionByID(createdEntityID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(entityFoundedInStorage.Cities) != 1 {
		t.Fatalf("Expect 1 city with id: %v but got: %v count", cityID, entityFoundedInStorage.Cities)
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

func createCity(entityForCreate storage.City, language string, databaseClient *dataBaseClient.Dgraph) (entityID string, err error) {

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

	err = addOtherLanguageForCityName(uid, entityForCreate.Name, language, databaseClient)
	if err != nil {
		return uid, nil
	}

	return uid, nil
}

func addOtherLanguageForCityName(cityID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	namePredicate := fmt.Sprintf(`<%s> <cityName> %s .`, cityID, "\""+name+"\""+"@"+language)

	mutation := &dataBaseAPI.Mutation{
		SetNquads: []byte(namePredicate),
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	_, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return err
	}

	return nil
}

func createInstruction(entityForCreate storage.Instruction, databaseClient *dataBaseClient.Dgraph) (string, error) {
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

func readCityByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.City, err error) {
	variables := struct {
		CityID   string
		Language string
	}{
		CityID:   entityID,
		Language: language}

	queryTemplate, err := template.New("ReadCityByID").Parse(`{
				cities(func: uid("{{.CityID}}")) @filter(has(cityName)) {
					uid
					cityName: cityName@{{.Language}}
					cityIsActive
				}
			}`)

	entity.ID = entityID

	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	response, err := query(queryBuf.String(), databaseClient)
	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	type entitiesInStore struct {
		Entities []storage.City `json:"cities"`
	}

	var foundedCategories entitiesInStore

	err = json.Unmarshal(response, &foundedCategories)
	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	if len(foundedCategories.Entities) == 0 {
		return entity, err
	}

	return foundedCategories.Entities[0], nil
}

func readInstructionByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.Instruction, err error) {
	variables := struct {
		InstructionID string
		Language      string
	}{
		InstructionID: entityID,
		Language:      language}

	queryTemplate, err := template.New("ReadEntityByID").Parse(`{
				instructions(func: uid("{{.InstructionID}}")) @filter(eq(instructionLanguage, {{.Language}})) {
					uid
					instructionLanguage
					instructionIsActive
					has_page {
						uid
						path
						pageInPaginationSelector
						pageParamPath
						cityParamPath
						itemSelector
						nameOfItemSelector
						priceOfItemSelector
					}
					has_city @filter(eq(cityIsActive, true)) {
						uid
						cityName: cityName@{{.Language}}
						cityIsActive
					}
					has_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIsActive
					}
					has_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
					}
				}
			}`)

	entity = storage.Instruction{ID: entityID}

	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	response, err := query(queryBuf.String(), databaseClient)
	if err != nil {
		ExecutorLogger.Println(err)
	}

	type EntitiesInStorage struct {
		Entities []storage.Instruction `json:"instructions"`
	}

	var foundedEntities EntitiesInStorage

	err = json.Unmarshal(response, &foundedEntities)
	if err != nil {
		ExecutorLogger.Println(err)
	}

	if len(foundedEntities.Entities) == 0 {
		return entity, errors.New("Entity does not exist")
	}

	return foundedEntities.Entities[0], nil
}

func query(request string, client *dataBaseClient.Dgraph) (response []byte, err error) {
	transaction := client.NewTxn()
	resp, err := transaction.Query(context.Background(), request)
	if err != nil {
		return response, err
	}

	return resp.GetJson(), nil
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
