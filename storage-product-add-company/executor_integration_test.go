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

func TestExecutor_AddCompanyToProduct(t *testing.T) {
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

	createdProductID, err := createProduct(entityForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdProductID == "" {
		t.Fatalf("Created entity id is empty")
	}

	defer func() {
		err = deleteEntityByID(createdProductID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	err = addOtherLanguageForProductName(createdProductID, NameOfProduct, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	productFoundedInStorage, err := readProductByID(createdProductID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if productFoundedInStorage.ID != createdProductID {
		t.Fatalf("ID: %v of founded entity in storage is not ID: %v of created entity", productFoundedInStorage.ID, createdProductID)
	}

	if productFoundedInStorage.Name != NameOfProduct {
		t.Fatalf("Name: %v of founded entity in storage is not Name: %v of created entity", productFoundedInStorage.Name, NameOfProduct)
	}

	if len(productFoundedInStorage.Companies) > 0 {
		t.Fatalf("Expect 0 companies but got: %v", productFoundedInStorage.Companies)
	}

	companyForCreate := storage.Company{
		Name:     "Test company name",
		IRI:      "//",
		IsActive: true}

	schema = `
		companyName: string @lang @index(term) .
		companyIsActive: bool @index(bool) .
		has_category: uid @count .
	`

	err = setUpSchema(schema, databaseClient)

	companyID, err := createCompany(companyForCreate, Language, databaseClient)
	if err != nil {
		t.Fatalf("Company does not create")
	}

	defer func() {
		err = deleteEntityByID(companyID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	company, err := readCompanyByID(companyID, Language, databaseClient)
	if err != nil {
		t.Fatalf("Company does not exist: %v", err)
	}

	if company.ID != companyID {
		t.Fatalf("ID: %v of founded company in storage is not ID: %v of created company", company.ID, companyID)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}

	err = executor.AddCompanyToProduct(createdProductID, companyID)
	if err != nil {
		t.Fatalf("Name of Company does not added")
	}

	productFoundedInStorage, err = readProductByID(createdProductID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(productFoundedInStorage.Companies) != 1 {
		t.Fatalf("Expect 1 company with id: %v but got: %v count", companyID, productFoundedInStorage.Companies)
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

func createCompany(entityForCreate storage.Company, language string, databaseClient *dataBaseClient.Dgraph) (entityID string, err error) {

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

	err = addOtherLanguageForCompanyName(uid, entityForCreate.Name, language, databaseClient)
	if err != nil {
		return uid, nil
	}

	return uid, nil
}

func addOtherLanguageForCompanyName(companyID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	namePredicate := fmt.Sprintf(`<%s> <companyName> %s .`, companyID, "\""+name+"\""+"@"+language)

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

func createProduct(entityForCreate storage.Product, databaseClient *dataBaseClient.Dgraph) (string, error) {
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

func addOtherLanguageForProductName(productID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	namePredicate := fmt.Sprintf(`<%s> <productName> %s .`, productID, "\""+name+"\""+"@"+language)

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

func readCompanyByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.Company, err error) {
	variables := struct {
		CompanyID string
		Language  string
	}{
		CompanyID: entityID,
		Language:  language}

	queryTemplate, err := template.New("ReadCompanyByID").Parse(`{
				companies(func: uid("{{.CompanyID}}")) @filter(has(companyName)) {
					uid
					companyName: companyName@{{.Language}}
					companyIri
					companyIsActive
					has_category @filter(eq(categoryIsActive, true)) {
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
							}
						}
						has_product @filter(uid_in(belongs_to_company, {{.CompanyID}}) AND eq(productIsActive, true)) {
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
							has_price @filter(eq(priceIsActive, true)) {
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

	type entitiesInStore struct {
		Entities []storage.Company `json:"companies"`
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

func readProductByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.Product, err error) {
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

	response, err := query(queryBuf.String(), databaseClient)
	if err != nil {
		return entity, err
	}

	type productsInStore struct {
		Products []storage.Product `json:"products"`
	}

	var foundedProducts productsInStore

	err = json.Unmarshal(response, &foundedProducts)
	if err != nil {
		return entity, err
	}

	if len(foundedProducts.Products) == 0 {
		return entity, err
	}

	return foundedProducts.Products[0], nil
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
