package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language        string
	CompanyID       string
	DatabaseGateway string
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
	company, err := executor.ReadCategoryByID(request.CompanyID, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadCategoryByID error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "ReadCategoryByID error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

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

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal category error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return string(encodedCategory)
}
