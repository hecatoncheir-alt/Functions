package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language        string
	InstructionID   string
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
	instruction, err := executor.ReadInstructionByID(request.InstructionID, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadInstructionByID error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "ReadInstructionByID error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedInstruction, err := json.Marshal(instruction)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal instruction error: %v. Error: %v", instruction, err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "Unmarshal instruction error",
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
