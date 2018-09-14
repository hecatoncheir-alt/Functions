package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language        string
	CategoryName    string
	DatabaseGateway string
}

type Response struct{ Message, Data, Error string }

// Handle a serverless request
func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)

		fmt.Println(warning)

		errorResponse := Response{Error: err.Error(), Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: request.DatabaseGateway}}
	categories, err := executor.ReadCategoriesByName(request.CategoryName, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadCategoriesByName error: %v", err)

		fmt.Println(warning)

		errorResponse := Response{Error: err.Error(), Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedCategories, err := json.Marshal(categories)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal categories error: %v. Error: %v", categories, err)

		fmt.Println(warning)

		errorResponse := Response{Error: err.Error(), Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	response := Response{Data: string(encodedCategories)}
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedResponse)
}
