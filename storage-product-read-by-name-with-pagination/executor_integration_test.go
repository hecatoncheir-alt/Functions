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

func TestExecutor_ReadProductsByNameWithPagination(t *testing.T) {
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

	// Product 0 with other language
	product0, err := createProduct(
		storage.Product{Name: "Первый тестовый продукт"}, "en", databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() {
		err = deleteEntityByID(product0.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Product 1
	product1, err := createProduct(
		storage.Product{Name: "Первый тестовый продукт"}, "ru", databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() {
		err = deleteEntityByID(product1.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Product 2
	product2, err := createProduct(
		storage.Product{Name: "Второй тестовый продукт"}, "ru", databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() {
		err = deleteEntityByID(product2.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Product 3
	product3, err := createProduct(
		storage.Product{Name: "Третий тестовый продукт"}, "ru", databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() {
		err = deleteEntityByID(product3.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Product 4
	product4, err := createProduct(
		storage.Product{Name: "Четвёртый тестовый продукт"}, "ru", databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() {
		err = deleteEntityByID(product4.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	// Product 5
	product5, err := createProduct(
		storage.Product{Name: "Пятый тестовый продукт"}, "ru", databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() {
		err = deleteEntityByID(product5.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}

	// First page
	foundedProductsForFirstPage, err := executor.ReadProductsByNameWithPagination("тестовый", "ru", 1, 2)
	if err != nil {
		t.Error(err)
	}
	if foundedProductsForFirstPage.TotalProductsFound != 5 {
		t.Errorf("Expected 5 products, actual: %v", foundedProductsForFirstPage.TotalProductsFound)
	}

	if foundedProductsForFirstPage.TotalProductsForOnePage != 2 {
		t.Errorf("Expected 2 products on one page, actual: %v", foundedProductsForFirstPage.TotalProductsForOnePage)
	}

	if foundedProductsForFirstPage.CurrentPage != 1 {
		t.Errorf("Expected page 1, actual: %v", foundedProductsForFirstPage.CurrentPage)
	}

	if foundedProductsForFirstPage.Products[0].Name != "Первый тестовый продукт" {
		t.Errorf("Expected \"Первый тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[0].Name)
	}

	if foundedProductsForFirstPage.Products[1].Name != "Второй тестовый продукт" {
		t.Errorf("Expected \"Второй тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[1].Name)
	}

	// Second page
	foundedProductsForSecondPage, err := executor.ReadProductsByNameWithPagination("тестовый", "ru", 2, 2)
	if err != nil {
		t.Error(err)
	}

	if foundedProductsForSecondPage.Products[0].Name != "Третий тестовый продукт" {
		t.Errorf("Expected \"Третий тестовый продукт\", actual: %v", foundedProductsForSecondPage.Products[0].Name)
	}

	if foundedProductsForSecondPage.Products[1].Name != "Четвёртый тестовый продукт" {
		t.Errorf("Expected \"Четвёртый тестовый продукт\", actual: %v", foundedProductsForSecondPage.Products[1].Name)
	}

	// Third page
	foundedProductsForThirdPage, err := executor.ReadProductsByNameWithPagination("тестовый", "ru", 3, 2)
	if err != nil {
		t.Error(err)
	}

	if foundedProductsForThirdPage.CurrentPage != 3 {
		t.Errorf("Expected page 3, actual: %v", foundedProductsForFirstPage.CurrentPage)
	}

	if len(foundedProductsForThirdPage.Products) != 1 {
		t.Errorf("Expected 1 product for one page, actual: %v", len(foundedProductsForFirstPage.Products))
	}

	if len(foundedProductsForThirdPage.Products) != 1 {
		t.Errorf("Expected one product. actual: %v", len(foundedProductsForFirstPage.Products))
	}

	if foundedProductsForThirdPage.Products[0].Name != "Пятый тестовый продукт" {
		t.Errorf("Expected \"Пятый тестовый продукт\", actual: %v", foundedProductsForFirstPage.Products[0].Name)
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

func createProduct(product storage.Product, language string, databaseClient *dataBaseClient.Dgraph) (storage.Product, error) {
	product.IsActive = true

	productID, err := createEntity(product, databaseClient)
	if err != nil {
		return product, nil
	}

	err = addOtherLanguageForEntityName(productID, product.Name, language, databaseClient)
	if err != nil {
		return product, err
	}

	return readProductByID(productID, language, databaseClient)
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

func readProductByID(productID, language string, databaseClient *dataBaseClient.Dgraph) (storage.Product, error) {
	product := storage.Product{}

	variables := struct {
		ProductID string
		Language  string
	}{
		ProductID: productID,
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

	product.ID = productID

	if err != nil {
		return product, err
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		return product, err
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return product, ErrProductsByNameCanNotBeFound
	}

	type entitiesInStorage struct {
		AllEntitiesFoundedByID []storage.Product `json:"products"`
	}

	var foundedEntities entitiesInStorage

	err = json.Unmarshal(response.GetJson(), &foundedEntities)
	if err != nil {
		log.Println(err)
		return product, ErrProductsByNameCanNotBeFound
	}

	if len(foundedEntities.AllEntitiesFoundedByID) == 0 {
		return product, ErrProductsByNameCanNotBeFound
	}

	return foundedEntities.AllEntitiesFoundedByID[0], nil
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
