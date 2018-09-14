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
	"os"
	"testing"
	"text/template"
)

func TestExecutor_AddCategoryToProduct(t *testing.T) {
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

	Language := "ru"
	NameOfProduct := "Test product name"

	entityForCreate := storage.Product{
		Name:     NameOfProduct,
		IsActive: true}

	createdEntityID, err := createEntity(entityForCreate, databaseClient)
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

	err = addOtherLanguageForCompanyName(createdEntityID, NameOfProduct, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	entityFoundedInStorage, err := readEntityByID(createdEntityID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if entityFoundedInStorage.ID != createdEntityID {
		t.Fatalf("ID: %v of founded entity in storage is not ID: %v of created entity", entityFoundedInStorage.ID, createdEntityID)
	}

	if entityFoundedInStorage.Name != NameOfProduct {
		t.Fatalf("Name: %v of founded entity in storage is not Name: %v of created entity", entityFoundedInStorage.Name, NameOfProduct)
	}

	if len(entityFoundedInStorage.Categories) > 0 {
		t.Fatalf("Expect 0 categories but got: %v", entityFoundedInStorage.Categories)
	}

	schema = `
		categoryName: string @lang @index(term).
		categoryIsActive: bool @index(bool) .
		belongs_to_company: uid .
		has_product: uid .
	`

	err = setUpSchema(schema, databaseClient)

	categoryForCreate := storage.Category{
		Name:     "Test category name",
		IsActive: true}

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

	err = executor.AddCategoryToProduct(createdEntityID, categoryID)
	if err != nil {
		t.Fatalf("Name of Category does not added")
	}

	entityFoundedInStorage, err = readEntityByID(createdEntityID, Language, databaseClient)
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
		return entity, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		return entity, err
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		return entity, err
	}

	type categoriesInStore struct {
		Categories []storage.Category `json:"categories"`
	}

	var foundedCategories categoriesInStore

	err = json.Unmarshal(response.GetJson(), &foundedCategories)
	if err != nil {
		return entity, err
	}

	if len(foundedCategories.Categories) == 0 {
		return entity, err
	}

	return foundedCategories.Categories[0], nil
}

func readEntityByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.Product, err error) {
	variables := struct {
		ProductID string
		Language  string
	}{
		ProductID: entityID,
		Language:  language}

	queryTemplate, err := template.New("ReadProductByID").Parse(`{
				products(func: uid("{{.ProductID}}")) @filter(has(productName)) {
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
					has_price @filter(eq(priceIsActive, true)) (orderasc: priceDateTime) {
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
							has_price @filter(eq(priceIsActive, true)) {
								uid
								priceValue
								priceDateTime
								priceCity
								priceIsActive
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

	entity.ID = entityID

	if err != nil {
		return entity, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		return entity, err
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		return entity, err
	}

	type productsInStore struct {
		Products []storage.Product `json:"products"`
	}

	var foundedProducts productsInStore

	err = json.Unmarshal(response.GetJson(), &foundedProducts)
	if err != nil {
		return entity, err
	}

	if len(foundedProducts.Products) == 0 {
		return entity, err
	}

	return foundedProducts.Products[0], nil
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
