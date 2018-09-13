package function

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"strconv"
)

// Instructions is a structure of settings for pages count parse
type Instructions struct {
	PageInPaginationSelector string `json:"pageInPaginationSelector,omitempty"`
}

type Request struct {
	IRI          string
	Instructions Instructions
}

type Response struct{ Message, Data, Error string }

func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)
		fmt.Println(warning)
	}

	pagesCount, err := getPagesCount(request.IRI, request.Instructions)
	if err != nil {
		warning := fmt.Sprintf(
			"Get count of pages error by IRI: %v. Error: %v",
			request.IRI,
			err)

		fmt.Println(warning)

		response := Response{
			Message: warning,
			Error:   err.Error(),
		}

		encodedResponse, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}

		return string(encodedResponse)
	}

	encodedPagesCount, err := json.Marshal(pagesCount)
	if err != nil {
		fmt.Println(err)
	}

	response := Response{Data: string(encodedPagesCount)}

	encodedResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}

	return string(encodedResponse)
}

func getPagesCount(pageIRI string, instructions Instructions) (pagesCount int, err error) {
	collector := colly.NewCollector(colly.Async(true))

	collector.OnHTML(instructions.PageInPaginationSelector,
		func(element *colly.HTMLElement) {
			pagesCount, err = strconv.Atoi(element.Text)

			if err != nil {
				warning := fmt.Sprintf(
					"Get count of pages from: %v failed with response: %v. Error: %v",
					element.Request.URL,
					element.Response.Body,
					err)

				fmt.Println(warning)
			}
		})

	collector.OnError(func(response *colly.Response, err error) {
		warning := fmt.Sprintf(
			"Request URL: %v failed with response: %v. Error: %v",
			response.Request.URL,
			response,
			err)

		fmt.Println(warning)
	})

	err = collector.Visit(pageIRI)
	if err != nil {
		warning := fmt.Sprintf(
			"Error visit URL: %v. Error: %v",
			pageIRI,
			err)

		fmt.Println(warning)

		return 0, err
	}

	collector.Wait()

	return pagesCount, nil
}
