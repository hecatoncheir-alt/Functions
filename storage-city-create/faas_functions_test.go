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

func TestFAASFunctions_ReadCitiesByName(t *testing.T) {

	LanguageForTest := "ru"
	CityNameForTest := "TestCityName"
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

		if responseBodyEncoded["CityName"] != CityNameForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", CityNameForTest, responseBodyEncoded["CityName"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedCitiesInStorage := []storage.City{
			{
				ID:       "0x12",
				Name:     "Test city name",
				IsActive: true},
			{
				ID:       "0x13",
				Name:     "Other test city name",
				IsActive: true}}

		encodedExistedCitiesInStorage, err := json.Marshal(existedCitiesInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		response := Response{Data: string(encodedExistedCitiesInStorage)}

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
	cities := faas.ReadCitiesByName(CityNameForTest, LanguageForTest)

	if len(cities) < 1 {
		t.Fatalf("Expect more cities that 1, but got: %v", len(cities))
	}
}

func TestFAASFunctions_ReadCompanyByID(t *testing.T) {
	LanguageForTest := "ru"
	CityIDForTest := "0x12"
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

		if responseBodyEncoded["CityID"] != CityIDForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", CityIDForTest, responseBodyEncoded["CityID"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedCityInStorage := storage.City{
			ID:       "0x12",
			Name:     "Test city name",
			IsActive: true}

		encodedExistedCityInStorage, err := json.Marshal(existedCityInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		response := Response{Data: string(encodedExistedCityInStorage)}

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
	city := faas.ReadCityByID(CityIDForTest, LanguageForTest)

	if city.ID != "0x12" {
		t.Fatalf("Expect city id: %v, but got: %v", CityIDForTest, city.ID)
	}
}
