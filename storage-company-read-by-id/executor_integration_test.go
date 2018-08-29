package function

import (
	"context"
	"encoding/json"
	"fmt"
	dataBaseClient "github.com/dgraph-io/dgo"
	dataBaseAPI "github.com/dgraph-io/dgo/protos/api"
	"github.com/hecatoncheir/Storage"
	"google.golang.org/grpc"
	"testing"
)

func TestExecutor_ReadCompanyByID(t *testing.T) {
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
	Language := "ru"

	executor := Executor{Store: &storage.Store{DatabaseGateway: DatabaseGateway}}
	_, err = executor.ReadCompanyByID(CompanyID, Language)
	if err != ErrCompanyCanNotBeWithoutID {
		t.Error(err)
	}

	FakeCompanyID := "0x12"
	_, err = executor.ReadCompanyByID(FakeCompanyID, Language)
	if err != ErrCompanyDoesNotExist {
		t.Error(err)
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

	err = addOtherLanguageForCompanyName(createdCompanyID, CompanyName, Language, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	companyFoundedInStorage, err := executor.ReadCompanyByID(createdCompanyID, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if companyFoundedInStorage.ID != createdCompanyID {
		t.Fatalf("ID: %v of founded company in storage is not ID: %v of created company", companyFoundedInStorage.ID, createdCompanyID)
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
	encodedCity, err := json.Marshal(companyForCreate)
	if err != nil {
		return "", err
	}

	mutation := &dataBaseAPI.Mutation{
		SetJson:   encodedCity,
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
