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

func TestFAASFunctions_CompaniesReadByName(t *testing.T) {

	LanguageForTest := "ru"
	CompanyNameForTest := "TestCompanyName"
	DatabaseGatewayForTest := "http://TestDatabaseGateway"

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

		if responseBodyEncoded["CompanyName"] != CompanyNameForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", CompanyNameForTest, responseBodyEncoded["Language"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedCompaniesInStorage := []storage.Company{
			{
				ID:       "0x12",
				IRI:      "/",
				Name:     "Test company name",
				IsActive: true},
			{
				ID:       "0x13",
				IRI:      "/",
				Name:     "Other test company name",
				IsActive: true}}

		encodedExistedCompaniesInStorage, err := json.Marshal(existedCompaniesInStorage)

		io.WriteString(w, string(encodedExistedCompaniesInStorage))
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := FAASFunctions{FAASGateway: testServer.URL}
	companies := faas.CompaniesReadByName(CompanyNameForTest, LanguageForTest, DatabaseGatewayForTest)

	if len(companies) < 1 {
		t.Fatalf("Expect more companies that 1, but got: %v", len(companies))
	}
}
