package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParserCanParsePage(t *testing.T) {
	testFileContent, err := ioutil.ReadFile("handler_test_page.html")
	if err != nil {
		t.Errorf(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, err := w.Write(testFileContent)
		if err != nil {
			t.Errorf(err.Error())
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	request := Request{
		IRI: fmt.Sprint(server.URL, "/test"),
		Instructions: Instructions{
			ItemSelector:               ".c-product-tile",
			PreviewImageOfItemSelector: ".c-product-tile-picture__link .lazy-load-image-holder img",
			NameOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
			LinkOfItemSelector:         ".c-product-tile__description .sel-product-tile-title",
			PriceOfItemSelector:        ".c-product-tile__checkout-section .c-pdp-price__current"},
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		t.Errorf(err.Error())
	}

	encodedResponse := Handle(bytes)

	response := Response{}

	err = json.Unmarshal([]byte(encodedResponse), &response)
	if err != nil {
		t.Errorf(err.Error())
	}

	if response.Error != "" {
		t.Errorf(response.Error)
	}

	var listOfProducts []Product

	json.Unmarshal([]byte(response.Data), &listOfProducts)

	expectedLengthOfProductsList := 12

	if len(listOfProducts) != expectedLengthOfProductsList {
		t.Errorf("expected '%d' but got '%d'", expectedLengthOfProductsList, len(listOfProducts))
	}
}
