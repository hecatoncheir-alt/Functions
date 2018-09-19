package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct{ Language, ProductID, DatabaseGateway string }
type Response struct{ Message, Data, Error string }

// Handle a serverless request
func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)

		fmt.Println(warning)

		errorResponse := Response{Message: warning, Data: string(req), Error: err.Error()}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: request.DatabaseGateway}}
	product, err := executor.ReadProductByID(request.ProductID, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadProductByID error: %v", err)

		fmt.Println(warning)

		errorResponse := Response{Message: warning, Data: string(req), Error: err.Error()}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedProduct, err := json.Marshal(product)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal product error: %v. Error: %v", product, err)

		fmt.Println(warning)

		errorResponse := Response{Message: warning, Data: string(req), Error: err.Error()}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	response := Response{Data: string(encodedProduct)}

	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedResponse)
}
