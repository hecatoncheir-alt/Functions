package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParserCanParsePagesCount(t *testing.T) {
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
			PageInPaginationSelector: ".c-pagination > .c-pagination__num"},
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		t.Errorf(err.Error())
	}

	pagesCount := Handle(bytes)

	expectedPagesCount := 68

	if pagesCount != expectedPagesCount {
		t.Errorf("expected '%d' but got '%d'", expectedPagesCount, pagesCount)
	}
}
