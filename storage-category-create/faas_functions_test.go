package function

import (
	"encoding/json"
	"github.com/hecatoncheir/Storage"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFAASFunctions_ReadCategoriesByName(t *testing.T) {

	LanguageForTest := "ru"
	CategoryNameForTest := "Test category name"
	DatabaseGatewayForTest := "http://"

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Read body of request error: %v", err)
		}

		var responseBodyEncoded map[string]string
		err = json.Unmarshal(encodedBody, &responseBodyEncoded)
		if err != nil {
			t.Errorf("Unmarshal body of request error: %v", err)
		}

		if responseBodyEncoded["Language"] != LanguageForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", LanguageForTest, responseBodyEncoded["Language"])
		}

		if responseBodyEncoded["CategoryName"] != CategoryNameForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", CategoryNameForTest, responseBodyEncoded["CategoryName"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedCategoriesInStorage := []storage.Category{
			{
				ID:       "0x12",
				Name:     "Test category name",
				IsActive: true},
			{
				ID:       "0x13",
				Name:     "Other test category name",
				IsActive: true}}

		encodedExistedCategoriesInStorage, err := json.Marshal(existedCategoriesInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		response := Response{Data: string(encodedExistedCategoriesInStorage)}

		encodedResponse, err := json.Marshal(response)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = io.WriteString(w, string(encodedResponse))
		if err != nil {
			t.Error(err.Error())
		}
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := &FAASFunctions{FunctionsGateway: testServer.URL, DatabaseGateway: DatabaseGatewayForTest}
	categories := faas.ReadCategoriesByName(CategoryNameForTest, LanguageForTest)

	if len(categories) < 1 {
		t.Fatalf("Expect more categories that 1, but got: %v", len(categories))
	}
}

func TestFAASFunctions_ReadCategoryByID(t *testing.T) {

	LanguageForTest := "ru"
	CategoryIDForTest := "0x12"
	DatabaseGatewayForTest := "http://"

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("Read body of request error: %v", err)
		}

		var responseBodyEncoded map[string]string
		err = json.Unmarshal(encodedBody, &responseBodyEncoded)
		if err != nil {
			t.Errorf("Unmarshal body of request error: %v", err)
		}

		if responseBodyEncoded["Language"] != LanguageForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", LanguageForTest, responseBodyEncoded["Language"])
		}

		if responseBodyEncoded["CategoryID"] != CategoryIDForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", CategoryIDForTest, responseBodyEncoded["CategoryID"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedCategoriesInStorage := storage.Category{
			ID:       "0x12",
			Name:     "Test category name",
			IsActive: true}

		encodedExistedCategoriesInStorage, err := json.Marshal(existedCategoriesInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = io.WriteString(w, string(encodedExistedCategoriesInStorage))
		if err != nil {
			t.Error(err.Error())
		}
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := &FAASFunctions{FunctionsGateway: testServer.URL, DatabaseGateway: DatabaseGatewayForTest}
	category := faas.ReadCategoryByID(CategoryIDForTest, LanguageForTest)

	if category.ID != "0x12" {
		t.Fatalf("Expect category id %v, but got: %v", CategoryIDForTest, category.ID)
	}

	if category.Name == "" {
		t.Fatalf("Expect category name is not empty, but got: %v", category.Name)
	}

}
