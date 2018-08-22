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
	PageInstruction storage.PageInstruction
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

	executor := Executor{
		Store: &storage.Store{DatabaseGateway: request.DatabaseGateway},
		Functions: &FAASFunctions{
			DatabaseGateway:  request.DatabaseGateway,
			FunctionsGateway: request.FunctionsGateway}}

	createdPageInstruction, err := executor.CreatePageInstruction(request.PageInstruction, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"CreatePageInstruction error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "CreatePageInstruction error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedInstruction, err := json.Marshal(createdPageInstruction)
	if err != nil {
		warning := fmt.Sprintf(
			"Marshal page instruction error: %v. Error: %v", createdPageInstruction, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Marshal created page instruction error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return string(encodedInstruction)
}
