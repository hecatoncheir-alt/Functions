package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language,
	DatabaseGateway,
	FunctionsGateway string
	Category storage.Category
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

		errorResponse := Response{Error: err.Error(), Message: warning, Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	executor := Executor{
		Store: &storage.Store{DatabaseGateway: request.DatabaseGateway},
		Functions: &FAASFunctions{
			DatabaseGateway:  request.DatabaseGateway,
			FunctionsGateway: request.FunctionsGateway}}

	createdCompany, err := executor.CreateCategory(request.Category, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"CreateCategory error: %v", err)

		fmt.Println(warning)

		errorResponse := Response{Error: err.Error(), Message: warning, Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedCompany, err := json.Marshal(createdCompany)
	if err != nil {
		warning := fmt.Sprintf(
			"Marshal category error: %v. Error: %v", createdCompany, err)

		fmt.Println(warning)

		errorResponse := Response{Error: err.Error(), Message: warning, Data: string(req)}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	response := Response{Data: string(encodedCompany)}
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedResponse)
}
