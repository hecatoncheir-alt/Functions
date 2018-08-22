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

func TestFAASFunctions_ReadPageInstructionByID(t *testing.T) {
	LanguageForTest := "ru"
	PageInstructionIDForTest := "0x12"
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

		if responseBodyEncoded["PageInstructionID"] != PageInstructionIDForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", PageInstructionIDForTest, responseBodyEncoded["PageInstructionID"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedPageInstructionInStorage := storage.PageInstruction{
			ID:   "0x12",
			Path: "http://"}

		encodedExistedPageInstructionInStorage, err := json.Marshal(existedPageInstructionInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = io.WriteString(w, string(encodedExistedPageInstructionInStorage))
		if err != nil {
			t.Error(err.Error())
		}
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := &FAASFunctions{FunctionsGateway: testServer.URL, DatabaseGateway: DatabaseGatewayForTest}
	instruction := faas.ReadPageInstructionByID(PageInstructionIDForTest, LanguageForTest)

	if instruction.ID != "0x12" {
		t.Fatalf("Expect page instruction id: %v, but got: %v", PageInstructionIDForTest, instruction.ID)
	}
}
