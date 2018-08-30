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

func TestExecutor_AddCategoryToInstruction(t *testing.T) {
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

	if len(entityFoundedInStorage.Categories) > 0 {
		t.Fatalf("Expect 0 categories but got: %v", entityFoundedInStorage.Categories)
	}

	categoryForCreate := storage.Category{
		Name:     "Test category name",
		IsActive: true}

	schema = `
		categoryName: string @lang @index(term).
		categoryIsActive: bool @index(bool) .
		belongs_to_company: uid .
		has_product: uid .
	`

	err = setUpSchema(schema, databaseClient)

	categoryID, err := createCategory(categoryForCreate, Language, databaseClient)
	if err != nil {
		t.Fatalf("Category does not create")
	}

	defer func() {
		err = deleteEntityByID(categoryID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	category, err := readCategoryByID(categoryID, Language, databaseClient)
	if err != nil {
		t.Fatalf("Category does not exist: %v", err)
	}

	if category.ID != categoryID {
		t.Fatalf("ID: %v of founded category in storage is not ID: %v of created category", category.ID, categoryID)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}

	err = executor.AddCategoryToInstruction(createdEntityID, categoryID)
	if err != nil {
		t.Fatalf("Name of Category does not added")
	}

	entityFoundedInStorage, err = readInstructionByID(createdEntityID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(entityFoundedInStorage.Categories) != 1 {
		t.Fatalf("Expect 1 category with id: %v but got: %v count", categoryID, entityFoundedInStorage.Categories)
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

func createCategory(entityForCreate storage.Category, language string, databaseClient *dataBaseClient.Dgraph) (entityID string, err error) {

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

	err = addOtherLanguageForCategoryName(uid, entityForCreate.Name, language, databaseClient)
	if err != nil {
		return uid, nil
	}

	return uid, nil
}

func addOtherLanguageForCategoryName(categoryID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	forCategoryNamePredicate := fmt.Sprintf(`<%s> <categoryName> %s .`, categoryID, "\""+name+"\""+"@"+language)

	mutation := &dataBaseAPI.Mutation{
		SetNquads: []byte(forCategoryNamePredicate),
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

func readCategoryByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.Category, err error) {
	variables := struct {
		CategoryID string
		Language   string
	}{
		CategoryID: entityID,
		Language:   language}

	queryTemplate, err := template.New("ReadCategoryByID").Parse(`{
				categories(func: uid("{{.CategoryID}}")) @filter(has(categoryName)) {
					uid
					categoryName: categoryName@{{.Language}}
					categoryIsActive
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIsActive
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
							belong_to_company @filter(eq(companyIsActive, true)) {
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_product @filter(eq(productIsActive, true)) {
						uid
						productName: productName@{{.Language}}
						productIri
						previewImageLink
						productIsActive
						belongs_to_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
						}
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIsActive
						}
					}
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

	type categoriesInStore struct {
		Categories []storage.Category `json:"categories"`
	}

	var foundedCategories categoriesInStore

	err = json.Unmarshal(response, &foundedCategories)
	if err != nil {
		ExecutorLogger.Println(err)
		return entity, err
	}

	if len(foundedCategories.Categories) == 0 {
		return entity, err
	}

	return foundedCategories.Categories[0], nil
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
