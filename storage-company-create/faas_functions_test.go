package function

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"encoding/json"
)

func TestFAASFunctions_CompaniesReadByName(t *testing.T) {

	responseBody := `
		[
			{
				 "uid": "0x12"
				 "companyName":"Test company",
				 "companyIri":"/",
				 "companyIsActive":true
		  	},
			{  
				 "uid":"0x13",
				 "companyName":"Other test company",
				 "companyIri":"/",
				 "companyIsActive":true
			}
		]
	`

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

		if responseBodyEncoded[""] != ""{

		}



		io.WriteString(w, responseBody)
	})

	testServer := httptest.NewServer(testHandler)
	defer testServer.Close()

	faas := FAASFunctions{FAASGateway: testServer.URL}
	companies := faas.CompaniesReadByName("TestCompanyName", "ru", "")
	if len(companies) < 1 {
		t.Fatalf("Expect more companies that 1, but got: %v", len(companies))
	}
}
