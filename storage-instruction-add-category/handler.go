package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	DatabaseGateway,
	InstructionID,
	CategoryID string
}

type ErrorResponse struct {
	Error string
	Data  ErrorData
}

type ErrorData struct {
	Error   string
	Request string
}

type NoErrorResponse struct {
	Error string
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

	err = executor.AddCategoryToInstruction(request.InstructionID, request.CategoryID)
	if err != nil {
		warning := fmt.Sprintf(
			"Add Category to Instruction error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Add Category to instruction error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	noErrorResponse := NoErrorResponse{Error: ""}

	encodedNoErrorResponse, err := json.Marshal(noErrorResponse)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedNoErrorResponse)
}
