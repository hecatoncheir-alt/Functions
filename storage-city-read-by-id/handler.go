package function

import (
	"encoding/json"
	"fmt"
	"github.com/hecatoncheir/Storage"
)

type Request struct{ Language, CityID, DatabaseGateway string }
type Response struct{ Message, Error, Data string }

// Handle a serverless request
func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)

		fmt.Println(warning)

		errorResponse := Response{Data: string(req), Error: err.Error()}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	executor := Executor{Store: &storage.Store{DatabaseGateway: request.DatabaseGateway}}
	city, err := executor.ReadCityByID(request.CityID, request.Language)
	if err != nil {
		warning := fmt.Sprintf(
			"ReadCityByID error: %v", err)

		fmt.Println(warning)

		errorResponse := Response{Data: string(req), Error: err.Error()}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	encodedCity, err := json.Marshal(city)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal city error: %v. Error: %v", city, err)

		fmt.Println(warning)

		errorResponse := Response{Data: string(req), Error: err.Error()}
		response, err := json.Marshal(errorResponse)
		if err != nil {
			fmt.Println(err)
		}

		return string(response)
	}

	response := Response{Data: string(encodedCity)}
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedResponse)
}
