package function

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"log"
	"os"
	"testing"
	"text/template"
	"time"
)

func TestPricesOfProductMustBeSortedByDate(t *testing.T) {
	t.Skip("Database must be started")

	DatabaseGateway := os.Getenv("DatabaseGateway")
	if DatabaseGateway == "" {
		DatabaseGateway = "localhost:9080"
	}

	databaseClient, err := connectToDatabase(DatabaseGateway)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Product
	schema := `
		productName: string @lang @index(term, trigram) .
		productIri: string @index(term) .
		productImageLink: string @index(term) .
		productIsActive: bool @index(bool) .
		belongs_to_category: uid .
		belongs_to_company: uid .
	`

	err = setUpSchema(schema, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	createdProduct, err := createProduct(
		storage.Product{Name: "Первый тестовый продукт"}, "ru", databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() {
		err = deleteEntityByID(createdProduct.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Price
	schema = `
		pricesValue: float @index(float) .
		priceDateTime: dateTime @index(day) .
		priceIsActive: bool @index(bool) .
		belongs_to_city: uid @count .
		belongs_to_product: uid @count .
		belongs_to_company: uid @count .
	`

	err = setUpSchema(schema, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// First price
	exampleFirstPriceDateTime := "2017-05-01T16:27:18.543653798Z"
	firstDateTime, err := time.Parse(time.RFC3339, exampleFirstPriceDateTime)
	if err != nil {
		t.Error(err)
	}

	createdFirstPrice, err := createPrice(storage.Price{Value: 0.1, DateTime: firstDateTime}, "ru", databaseClient)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = deleteEntityByID(createdFirstPrice.ID, databaseClient)
		if err != nil {
			t.Error(err)
		}
	}()

	err = addPriceToProduct(createdProduct.ID, createdFirstPrice.ID, databaseClient)
	if err != nil {
		t.Error(err)
	}

	// Second price
	exampleSecondDateTime := "2017-06-01T16:27:18.543653798Z"
	secondDateTime, err := time.Parse(time.RFC3339, exampleSecondDateTime)
	if err != nil {
		t.Error(err)
	}

	createdSecondPrice, err := createPrice(storage.Price{Value: 0.2, DateTime: secondDateTime}, "ru", databaseClient)
	if err != nil {
		t.Error(err)
	}

	defer func() {
		err = deleteEntityByID(createdFirstPrice.ID, databaseClient)
		if err != nil {
			t.Error(err)
		}
	}()

	err = addPriceToProduct(createdProduct.ID, createdSecondPrice.ID, databaseClient)
	if err != nil {
		t.Error(err)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}

	foundedProductsForFirstPage, err := executor.ReadProductsByNameWithPagination("тестовый", "ru", 1, 1)
	if err != nil {
		t.Error(err)
	}

	if len(foundedProductsForFirstPage.Products) != 1 {
		t.Fatalf(err.Error())
	}

	if len(foundedProductsForFirstPage.Products[0].Prices) != 2 {
		t.Fatalf(err.Error())
	}

	if foundedProductsForFirstPage.Products[0].Prices[0].Value != 0.2 {
		t.Fatalf(err.Error())
	}
}

func createPrice(price storage.Price, language string, databaseClient *dataBaseClient.Dgraph) (storage.Price, error) {
	price.IsActive = true

	priceID, err := createPriceEntity(price, databaseClient)
	if err != nil {
		return price, nil
	}

	return readPriceByID(priceID, language, databaseClient)
}

func createPriceEntity(price storage.Price, databaseClient *dataBaseClient.Dgraph) (string, error) {
	price.IsActive = true

	encodedProduct, err := json.Marshal(price)
	if err != nil {
		return "", err
	}

	mutation := dataBaseAPI.Mutation{
		SetJson:   encodedProduct,
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	assigned, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return "", err
	}

	uidOfCreatedPrice := assigned.Uids["blank-0"]

	if err != nil {
		return "", err
	}

	return uidOfCreatedPrice, nil

}

func addPriceToProduct(productID, priceID string, databaseClient *dataBaseClient.Dgraph) error {
	err := addEntityToOtherEntity(priceID, "belongs_to_product", productID, databaseClient)
	if err != nil {
		return err
	}

	err = addEntityToOtherEntity(productID, "has_price", priceID, databaseClient)
	if err != nil {
		return err
	}

	return nil
}

func addEntityToOtherEntity(entityID, field, addedEntityID string, databaseClient *dataBaseClient.Dgraph) error {

	subject := entityID
	predicate := field
	object := addedEntityID

	final := fmt.Sprintf(`<%s> <%s> <%s> .`, subject, predicate, object)

	mutation := dataBaseAPI.Mutation{
		SetNquads: []byte(final),
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	_, err := transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}

func readPriceByID(priceID, language string, databaseClient *dataBaseClient.Dgraph) (storage.Price, error) {

	variables := struct {
		PriceID  string
		Language string
	}{
		PriceID:  priceID,
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

	price := storage.Price{ID: priceID}
	if err != nil {
		return price, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return price, err
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		return price, err
	}

	type PricesInStore struct {
		Prices []storage.Price `json:"prices"`
	}

	var foundedPrices PricesInStore

	err = json.Unmarshal(response.GetJson(), &foundedPrices)
	if err != nil {
		return price, err
	}

	if len(foundedPrices.Prices) == 0 {
		return price, err
	}

	return foundedPrices.Prices[0], nil
}
