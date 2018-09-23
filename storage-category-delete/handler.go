package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct{ CategoryID, DatabaseGateway string }
type Response struct{ Message, Data, Error string }

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
	err = executor.DeleteCategoryByID(request.CategoryID)
	if err != nil {
		warning := fmt.Sprintf(
			"DeleteCategoryByID error: %v", err)

		errorResponse := Response{Error: err.Error(), Message: warning, Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	response := Response{Data: string(req)}
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedResponse)
}
