package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct {
	Language        string
	ProductID       string
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
	product, err := executor.ReadProductByID(request.ProductID, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadProductByID error: %v", err)

		fmt.Println(warning)

		errorResponse := ErrorResponse{
			Error: "ReadProductByID error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

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

		errorResponse := ErrorResponse{
			Error: "Unmarshal product error",
			Data: ErrorData{
				Request: string(req),
				Error:   err.Error()}}

		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	return string(encodedProduct)
}
