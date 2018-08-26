package function

import (
	"encoding/json"
	"github.com/hecatoncheir/Storage"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestFAASFunctions_ReadPriceByID(t *testing.T) {

	LanguageForTest := "ru"
	PriceIDForTest := "0x12"
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

		if responseBodyEncoded["PriceID"] != PriceIDForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", PriceIDForTest, responseBodyEncoded["PriceID"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedPriceInStorage := storage.Price{
			ID:       "0x12",
			Value:    0.1,
			DateTime: time.Now().UTC(),
			IsActive: true}

		encodedExistedPriceInStorage, err := json.Marshal(existedPriceInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = io.WriteString(w, string(encodedExistedPriceInStorage))
		if err != nil {
			t.Error(err.Error())
		}
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := &FAASFunctions{FunctionsGateway: testServer.URL, DatabaseGateway: DatabaseGatewayForTest}
	price := faas.ReadPriceByID(PriceIDForTest, LanguageForTest)

	if price.ID != "0x12" {
		t.Fatalf("Expect price id %v, but got: %v", PriceIDForTest, price.ID)
	}

	if price.Value != 0.1 {
		t.Fatalf("Expect price valut is 0.1, but got: %v", price.Value)
	}
}
