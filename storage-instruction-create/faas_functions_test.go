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

func TestFAASFunctions_ReadInstructionByID(t *testing.T) {
	LanguageForTest := "ru"
	InstructionIDForTest := "0x12"
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

		if responseBodyEncoded["InstructionID"] != InstructionIDForTest {
			t.Fatalf("Expected: \"%v\", but got: %v", InstructionIDForTest, responseBodyEncoded["InstructionID"])
		}

		if responseBodyEncoded["DatabaseGateway"] != DatabaseGatewayForTest {
			t.Fatalf(
				"Expected: \"%v\", but got: %v", DatabaseGatewayForTest, responseBodyEncoded["DatabaseGateway"])
		}

		existedInstructionInStorage := storage.Instruction{
			ID:       "0x12",
			IsActive: true}

		encodedExistedInstructionInStorage, err := json.Marshal(existedInstructionInStorage)
		if err != nil {
			t.Error(err.Error())
		}

		_, err = io.WriteString(w, string(encodedExistedInstructionInStorage))
		if err != nil {
			t.Error(err.Error())
		}
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := &FAASFunctions{FunctionsGateway: testServer.URL, DatabaseGateway: DatabaseGatewayForTest}
	instruction := faas.ReadInstructionByID(InstructionIDForTest, LanguageForTest)

	if instruction.ID != "0x12" {
		t.Fatalf("Expect instruction id: %v, but got: %v", InstructionIDForTest, instruction.ID)
	}
}
