package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language,
	ProductName,
	DatabaseGateway string

	CurrentPage,
	ItemsPerPage int
}

type ErrorResponse struct {
	Error string
	Data  ErrorData
}

type ErrorData struct {
	Error   string
	Request string
}

// Handle a serverless request
func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal request error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: request.DatabaseGateway}}
	companies, err := executor.ReadProductsByNameWithPagination(request.ProductName, request.Language, request.CurrentPage, request.ItemsPerPage)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadProductsByName error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "ReadProductsByName error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedProducts, err := json.Marshal(companies)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal products error: %v. Error: %v", companies, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal products error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return string(encodedProducts)
}
