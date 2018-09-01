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
	"os"
	"testing"
	"text/template"
)

func TestExecutor_ReadProductsByName(t *testing.T) {
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

	entityTestName := "Test product name"
	entityTestIRI := "//"

	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadProductsByName(entityTestName, Language)
	if err != ErrProductsByNameNotFound {
		t.Error(err)
	}

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

	err = addOtherLanguageForEntityName(createdEntityID, entityTestName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	entitiesByNameFromDatabase, err := readEntitiesByName(entityTestName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(entitiesByNameFromDatabase) < 1 {
		t.Fatalf("No one entity by name: %v found in database", entityTestName)
	}

	foundedEntities, err := executor.ReadProductsByName(entityTestName, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(foundedEntities) < 1 {
		t.Fatalf("No one entity by name: %v found in database", entityTestName)
	}

	if foundedEntities[0].Name != entityTestName {
		t.Fatalf("Name: %v of founded entity in storage is not Name: %v of created entity", foundedEntities[0].Name, entityTestName)
	}

	if foundedEntities[0].IRI != entityTestIRI {
		t.Fatalf("IRI: %v of founded entity in storage is not IRI: %v of created entity", foundedEntities[0].IRI, entityTestIRI)
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

func addOtherLanguageForEntityName(entityID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
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

func readEntitiesByName(entityName, language string, databaseClient *dataBaseClient.Dgraph) ([]storage.Product, error) {

	variables := struct {
		ProductName string
		Language    string
	}{
		ProductName: entityName,
		Language:    language}

	queryTemplate, err := template.New("ReadProductsByName").Parse(`{
				products(func: regexp(productName@{{.Language}}, /{{.ProductName}}/)) 
				@filter(eq(productIsActive, true) AND has(productName)) {
					uid
					productName: productName@{{.Language}}
					productIri
					previewImageLink
					productIsActive
					belongs_to_category @filter(eq(categoryIsActive, true)) {
						uid
						categoryName: categoryName@{{.Language}}
						categoryIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							has_category @filter(eq(categoryIsActive, true)) {
								uid
								categoryName: categoryName@{{.Language}}
								categoryIsActive
								belong_to_company @filter(eq(companyIsActive, true)){
									uid
									companyName: companyName@{{.Language}}
									companyIsActive
								}
							}
						}
					}
					belongs_to_company @filter(eq(companyIsActive, true)) {
						uid
						companyName: companyName@{{.Language}}
						companyIri
						companyIsActive
						has_category @filter(eq(categoryIsActive, true)) {
							uid
							categoryName: categoryName@{{.Language}}
							categoryIsActive
							belong_to_company @filter(eq(companyIsActive, true)){
								uid
								companyName: companyName@{{.Language}}
								companyIsActive
							}
						}
					}
					has_price @filter(eq(priceIsActive, true)) {
						uid
						priceValue
						priceDateTime
						priceIsActive
						belongs_to_company @filter(eq(companyIsActive, true)) {
							uid
							companyName: companyName@{{.Language}}
							companyIri
							companyIsActive
						}
						belongs_to_product @filter(eq(productIsActive, true)) {
							uid
							productName: productName@{{.Language}}
							productIri
							previewImageLink
							productIsActive
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceIsActive
								belongs_to_company @filter(eq(companyIsActive, true)) {
									uid
									companyName: companyName@{{.Language}}
									companyIri
									companyIsActive
								}
							}
						}
						belongs_to_city @filter(eq(cityIsActive, true)) {
							uid
							cityName: cityName@{{.Language}}
							cityIsActive
						}
					}
				}
			}`)

	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	type entitiesInStorage struct {
		AllEntitiesFoundedByName []storage.Product `json:"products"`
	}

	var foundedEntities entitiesInStorage

	err = json.Unmarshal(response.GetJson(), &foundedEntities)
	if err != nil {
		log.Println(err)
		return nil, ErrProductsByNameCanNotBeFound
	}

	if len(foundedEntities.AllEntitiesFoundedByName) == 0 {
		return nil, ErrProductsByNameCanNotBeFound
	}

	return foundedEntities.AllEntitiesFoundedByName, nil
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
