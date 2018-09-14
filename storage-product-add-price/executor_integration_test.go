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
	"time"
)

func TestExecutor_AddPriceToProduct(t *testing.T) {
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

	if len(productFoundedInStorage.Prices) > 0 {
		t.Fatalf("Expect 0 prices but got: %v", productFoundedInStorage.Prices)
	}

	schema = `
		pricesValue: float @index(float) .
		priceDateTime: dateTime @index(day) .
		priceIsActive: bool @index(bool) .
		belongs_to_city: uid @count .
		belongs_to_product: uid @count .
		belongs_to_company: uid @count .
	`

	err = setUpSchema(schema, databaseClient)

	exampleFirstPriceDateTime := "2017-05-01T16:27:18.543653798Z"
	firstDateTime, err := time.Parse(time.RFC3339, exampleFirstPriceDateTime)
	if err != nil {
		t.Error(err)
	}

	priceID, err := createPrice(
		storage.Price{Value: 0.1, DateTime: firstDateTime, IsActive: true}, databaseClient)
	if err != nil {
		t.Fatalf("Price does not create")
	}

	defer func() {
		err = deleteEntityByID(priceID, databaseClient)
		if err != nil {
			t.Error(err)
		}
	}()

	price, err := readPriceByID(priceID, Language, databaseClient)
	if err != nil {
		t.Fatalf("Price does not exist: %v", err)
	}

	if price.ID != priceID {
		t.Fatalf("ID: %v of founded price in storage is not ID: %v of created price", price.ID, priceID)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}

	err = executor.AddPriceToProduct(createdProductID, priceID)
	if err != nil {
		t.Fatalf("Name of Price does not added")
	}

	productFoundedInStorage, err = readProductByID(createdProductID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(productFoundedInStorage.Prices) != 1 {
		t.Fatalf("Expect 1 price with id: %v but got: %v count", priceID, productFoundedInStorage.Prices)
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

func createPrice(entityForCreate storage.Price, databaseClient *dataBaseClient.Dgraph) (entityID string, err error) {

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

func readPriceByID(entityID, language string, databaseClient *dataBaseClient.Dgraph) (entity storage.Price, err error) {

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

	type PricesInStore struct {
		Prices []storage.Price `json:"prices"`
	}

	var foundedPrices PricesInStore

	err = json.Unmarshal(response.GetJson(), &foundedPrices)
	if err != nil {
		return entity, err
	}

	if len(foundedPrices.Prices) == 0 {
		return entity, err
	}

	return foundedPrices.Prices[0], nil
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

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		return entity, err
	}
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
