package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct{ Language, CategoryID, DatabaseGateway string }
type Response struct{ Message, Error, Data string }

// Handle a serverless request
func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)

		errorResponse := Response{Error: err.Error(), Message: warning, Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: request.DatabaseGateway}}
	company, err := executor.ReadCategoryByID(request.CategoryID, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadCategoryByID error: %v", err)

		errorResponse := Response{Error: err.Error(), Message: warning, Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedCategory, err := json.Marshal(company)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal category error: %v. Error: %v", company, err)

		errorResponse := Response{Error: err.Error(), Message: warning, Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	response := Response{Data: string(encodedCategory)}
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedResponse)
}
