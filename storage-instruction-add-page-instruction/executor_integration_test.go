package function

import (
	"bytes"
	"context"
	"encoding/json"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"os"
	"testing"
	"text/template"
)

func TestExecutor_AddPageInstructionToInstruction(t *testing.T) {
	//t.Skip("Database must be started")

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

	if len(entityFoundedInStorage.PagesInstruction) > 0 {
		t.Fatalf("Expect 0 PagesInstruction but got: %v", entityFoundedInStorage.PagesInstruction)
	}

	pageInstructionForCreate := storage.PageInstruction{
		Path: "//"}

	schema = `
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

	pageInstructionID, err := createPageInstruction(pageInstructionForCreate, databaseClient)
	if err != nil {
		t.Fatalf("PageInstruction does not create")
	}

	defer func() {
		err = deleteEntityByID(pageInstructionID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	pageInstruction, err := readPageInstructionByID(pageInstructionID, Language, databaseClient)
	if err != nil {
		t.Fatalf("PageInstruction does not exist: %v", err)
	}

	if pageInstruction.ID != pageInstructionID {
		t.Fatalf("ID: %v of founded pageInstruction in storage is not ID: %v of created pageInstruction", pageInstruction.ID, pageInstructionID)
	}

	if pageInstruction.Path != "//" {
		t.Fatalf("Path: %v of founded pageInstruction in storage is not Path: %v of created pageInstruction", pageInstruction.ID, pageInstructionForCreate.Path)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}

	err = executor.AddPageInstructionToInstruction(createdEntityID, pageInstructionID)
	if err != nil {
		t.Fatalf("Name of PageInstruction does not added")
	}

	entityFoundedInStorage, err = readInstructionByID(createdEntityID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(entityFoundedInStorage.PagesInstruction) != 1 {
		t.Fatalf("Expect 1 pageInstruction with id: %v but got: %v count", pageInstructionID, entityFoundedInStorage.PagesInstruction)
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

func createPageInstruction(entityForCreate storage.PageInstruction, databaseClient *dataBaseClient.Dgraph) (entityID string, err error) {

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

func readPageInstructionByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.PageInstruction, err error) {

	variables := struct {
		PageInstructionID string
		Language          string
	}{
		PageInstructionID: entityID,
		Language:          language}

	queryTemplate, err := template.New("ReadPageInstructionByID").Parse(`{
				pageInstructions(func: uid("{{.PageInstructionID}}")) @filter(has(path)) {
					uid
					path
					pageInPaginationSelector
					pageParamPath
					cityParamPath
					itemSelector
					nameOfItemSelector
					priceOfItemSelector
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
		Entities []storage.PageInstruction `json:"pageInstructions"`
	}

	var foundedEntities entitiesInStore

	err = json.Unmarshal(response, &foundedEntities)
	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	if len(foundedEntities.Entities) == 0 {
		return entity, err
	}

	return foundedEntities.Entities[0], nil
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
