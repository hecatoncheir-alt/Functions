package function

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"log"
	"testing"
	"text/template"
)

func TestExecutor_DeleteCompanyByID(t *testing.T) {
	t.Skip("Database must be started")

	DatabaseGateway := "localhost:9080"
	databaseClient, err := connectToDatabase(DatabaseGateway)
	if err != nil {
		t.Fatalf(err.Error())
	}

	schema := `
		companyName: string @lang @index(term) .
		companyIsActive: bool @index(bool) .
		has_category: uid @count .
	`

	err = setUpCompanySchema(schema, databaseClient)

	CompanyID := ""

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	err = executor.DeleteCompanyByID(CompanyID)
	if err != ErrCompanyCanNotBeWithoutID {
		t.Fatalf(err.Error())
	}

	FakeCompanyID := "0x12"

	err = executor.DeleteCompanyByID(FakeCompanyID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	CompanyName := "Test company name"

	companyForCreate := storage.Company{
		Name:     CompanyName,
		IsActive: true}

	createdCompanyID, err := createCompany(companyForCreate, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCompanyID == "" {
		t.Fatalf("Created company id is empty")
	}

	Language := "ru"
	err = addOtherLanguageForCompanyName(createdCompanyID, CompanyName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	companyFromStore, err := readCompanyByID(createdCompanyID, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if companyFromStore.ID == "" {
		t.Fatalf("Created company not founded by id")
	}

	if companyFromStore.ID != createdCompanyID {
		t.Fatalf("Founded company id: %v not created company id: %v", companyFromStore.ID, createdCompanyID)
	}

	err = executor.DeleteCompanyByID(createdCompanyID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	companyFromStore, err = readCompanyByID(createdCompanyID, Language, databaseClient)
	if err.Error() != "company by id not found" {
		t.Fatalf(err.Error())
	}

	err = deleteCompanyByID(createdCompanyID, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	companyFromStore, err = readCompanyByID(createdCompanyID, Language, databaseClient)
	if err.Error() != "company by id not found" {
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

func setUpCompanySchema(schema string, databaseClient *dataBaseClient.Dgraph) error {
	operation := &dataBaseAPI.Operation{Schema: schema}

	err := databaseClient.Alter(context.Background(), operation)
	if err != nil {
		return err
	}

	return nil
}

func createCompany(companyForCreate storage.Company, databaseClient *dataBaseClient.Dgraph) (string, error) {
	encodedCompany, err := json.Marshal(companyForCreate)
	if err != nil {
		return "", err
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedCompany,
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	assigned, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return "", nil
	}

	uid := assigned.Uids["blank-0"]

	return uid, nil
}

func addOtherLanguageForCompanyName(companyID, name, language string, databaseClient *dataBaseClient.Dgraph) error {
	forCompanyNamePredicate := fmt.Sprintf(`<%s> <companyName> %s .`, companyID, "\""+name+"\""+"@"+language)

	mutation := &dataBaseAPI.Mutation{
		SetNquads: []byte(forCompanyNamePredicate),
		CommitNow: true}

	transaction := databaseClient.NewTxn()
	_, err := transaction.Mutate(context.Background(), mutation)
	if err != nil {
		return err
	}

	return nil
}

func readCompanyByID(companyID, language string, databaseClient *dataBaseClient.Dgraph) (storage.Company, error) {

	var (
		ErrCompanyByIDCanNotBeFound = errors.New("company by id can not be found")

		ErrCompanyDoesNotExist = errors.New("company by id not found")
	)

	variables := struct {
		CompanyID string
		Language  string
	}{
		CompanyID: companyID,
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

	company := storage.Company{ID: companyID}

	if err != nil {
		log.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	type companiesInStore struct {
		Companies []storage.Company `json:"companies"`
	}

	var foundedCompanies companiesInStore

	err = json.Unmarshal(response.GetJson(), &foundedCompanies)
	if err != nil {
		log.Println(err)
		return company, ErrCompanyByIDCanNotBeFound
	}

	if len(foundedCompanies.Companies) == 0 {
		return company, ErrCompanyDoesNotExist
	}

	return foundedCompanies.Companies[0], nil
}

func deleteCompanyByID(companyID string, databaseClient *dataBaseClient.Dgraph) error {
	deleteCompanyData, err := json.Marshal(map[string]string{"uid": companyID})
	if err != nil {
		return err
	}

	mutation := dataBaseAPI.Mutation{
		DeleteJson: deleteCompanyData,
		CommitNow:  true}

	transaction := databaseClient.NewTxn()

	_, err = transaction.Mutate(context.Background(), &mutation)
	if err != nil {
		return err
	}

	return nil
}
