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

func TestFAASFunctions_ReadProductsByName(t *testing.T) {

	LanguageForTest := "ru"
	ProductNameForTest := "Test product name"
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

		if responseBodyEncoded["ProductName"] != ProductNameForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", ProductNameForTest, responseBodyEncoded["ProductName"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedProductsInStorage := []storage.Product{
			{
				ID:       "0x12",
				Name:     "Test product name",
				IsActive: true},
			{
				ID:       "0x13",
				Name:     "Other test product name",
				IsActive: true}}

		encodedExistedProductsInStorage, err := json.Marshal(existedProductsInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = io.WriteString(w, string(encodedExistedProductsInStorage))
		if err != nil {
			t.Error(err.Error())
		}
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := &FAASFunctions{FunctionsGateway: testServer.URL, DatabaseGateway: DatabaseGatewayForTest}
	products := faas.ReadProductsByName(ProductNameForTest, LanguageForTest)

	if len(products) < 1 {
		t.Fatalf("Expect more products that 1, but got: %v", len(products))
	}
}

func TestFAASFunctions_ReadProductByID(t *testing.T) {

	LanguageForTest := "ru"
	ProductIDForTest := "0x12"
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

		if responseBodyEncoded["ProductID"] != ProductIDForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", ProductIDForTest, responseBodyEncoded["ProductID"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedProductInStorage := storage.Product{
			ID:       "0x12",
			Name:     "Test product name",
			IsActive: true}

		encodedExistedProductInStorage, err := json.Marshal(existedProductInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = io.WriteString(w, string(encodedExistedProductInStorage))
		if err != nil {
			t.Error(err.Error())
		}
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := &FAASFunctions{FunctionsGateway: testServer.URL, DatabaseGateway: DatabaseGatewayForTest}
	product := faas.ReadProductByID(ProductIDForTest, LanguageForTest)

	if product.ID != "0x12" {
		t.Fatalf("Expect product id %v, but got: %v", ProductIDForTest, product.ID)
	}

	if product.Name == "" {
		t.Fatalf("Expect product name is not empty, but got: %v", product.Name)
	}

}
