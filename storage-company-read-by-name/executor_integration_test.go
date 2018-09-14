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
	"testing"
	"text/template"
)

func TestExecutor_ReadCompaniesByName(t *testing.T) {
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

	CompanyName := "Test company name"
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadCompaniesByName(CompanyName, Language)
	if err != ErrCompaniesByNameNotFound {
		t.Error(err)
	}

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

	err = addOtherLanguageForCompanyName(createdCompanyID, CompanyName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	companyByNameFromDatabase, err := readCompaniesByName(CompanyName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(companyByNameFromDatabase) < 1 {
		t.Fatalf("No one company by name: %v found in database", CompanyName)
	}

	foundedCompanies, err := executor.ReadCompaniesByName(CompanyName, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(foundedCompanies) < 1 {
		t.Fatalf("No one company by name: %v found in database", CompanyName)
	}

	err = deleteCompanyByID(createdCompanyID, databaseClient)
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
	forCategoryNamePredicate := fmt.Sprintf(`<%s> <companyName> %s .`, companyID, "\""+name+"\""+"@"+language)

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

func readCompaniesByName(companyName, language string, databaseClient *dataBaseClient.Dgraph) ([]storage.Company, error) {

	variables := struct {
		CompanyName string
		Language    string
	}{
		CompanyName: companyName,
		Language:    language}

	queryTemplate, err := template.New("ReadCompaniesByName").Parse(`{
				companies(func: eq(companyName@{{.Language}}, "{{.CompanyName}}")) @filter(eq(companyIsActive, true)) {
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
						has_product @filter(eq(productIsActive, true)) { #TODO: belongs_to_company mast be an companyID
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

	if err != nil {
		log.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	queryBuf := bytes.Buffer{}
	err = queryTemplate.Execute(&queryBuf, variables)
	if err != nil {
		log.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	transaction := databaseClient.NewTxn()
	response, err := transaction.Query(context.Background(), queryBuf.String())
	if err != nil {
		log.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	type companiesInStorage struct {
		AllCompaniesFoundedByName []storage.Company `json:"companies"`
	}

	var foundedCompanies companiesInStorage
	err = json.Unmarshal(response.GetJson(), &foundedCompanies)
	if err != nil {
		log.Println(err)
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	if len(foundedCompanies.AllCompaniesFoundedByName) == 0 {
		return nil, ErrCompaniesByNameCanNotBeFound
	}

	return foundedCompanies.AllCompaniesFoundedByName, nil
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
