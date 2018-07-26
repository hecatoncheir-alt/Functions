package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Configuration"
)

type Request struct {
	CompanyName string
	APIVersion  string
	Language    string
}

type ErrorResponse struct {
	Error string
	APIVersion string
	Data  ErrorData
}

type ErrorData struct {
	Error   string
	Request string
}

// Handle a serverless request
func Handle(req []byte) string {
	config := configuration.New()

	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal request error",
			APIVersion: config.APIVersion,
			Data: ErrorData{
				Request: string(req),
				Error: err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return fmt.Sprintf("Hello, Go. You said: %s", string(req))
}
