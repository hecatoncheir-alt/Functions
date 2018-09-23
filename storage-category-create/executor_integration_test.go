package function

import (
	"github.com/hecatoncheir/Storage"
	"os"
	"testing"
	"time"
)

func TestIntegration_Executor(t *testing.T) {
	t.Skip("Database and FAAS must be started")

	DatabaseGateway := os.Getenv("DatabaseGateway")
	if DatabaseGateway == "" {
		DatabaseGateway = "192.168.99.101:31332"
	}

	FunctionsGateway := os.Getenv("FunctionsGateway")
	if FunctionsGateway == "" {
		FunctionsGateway = "http://192.168.99.101:31112/function"
	}

	databaseClient, err := connectToDatabase(DatabaseGateway)
	if err != nil {
		t.Fatalf(err.Error())
	}

	schema := `
		categoryName: string @lang @index(term).
		categoryIsActive: bool @index(bool) .
		belongs_to_company: uid .
		has_product: uid .
	`

	err = setUpSchema(schema, databaseClient)
	if err != nil {
		t.Fatalf(err.Error())
	}

	CategoryName := "Test category name"

	categoryForCreate := storage.Category{
		Name:     CategoryName,
		IsActive: true}

	executor := Executor{
		Store: &storage.Store{DatabaseGateway: DatabaseGateway},
		Functions: &FAASFunctions{
			DatabaseGateway:  DatabaseGateway,
			FunctionsGateway: FunctionsGateway}}

	Language := "ru"

	time.Sleep(3 * time.Second)

	createdCategory, err := executor.CreateCategory(categoryForCreate, Language)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if createdCategory.ID == "" {
		t.Fatalf("Created category id is empty")
	}

	defer func() {
		err = deleteCategoryByID(createdCategory.ID, databaseClient)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}()

	foundedCategory := executor.Functions.ReadCategoryByID(createdCategory.ID, Language)

	if foundedCategory.ID != createdCategory.ID {
		t.Fatalf("Created category: %v by id not found", createdCategory.ID)
	}

	if foundedCategory.Name != CategoryName {
		t.Fatalf("Created category: %v by name not found", CategoryName)
	}

	foundedCategories := executor.Functions.ReadCategoriesByName(foundedCategory.Name, Language)

	if len(foundedCategories) != 1 {
		t.Fatalf("No one category by name: %v found in database", CategoryName)
	}

	if foundedCategories[0].ID != createdCategory.ID {
		t.Fatalf("Created category: %v by id not found", createdCategory.ID)
	}

	if foundedCategories[0].Name != CategoryName {
		t.Fatalf("Created category: %v by name not found", CategoryName)
	}
}
