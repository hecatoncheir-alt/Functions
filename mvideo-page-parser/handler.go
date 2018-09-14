package function

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Instructions is a structure of settings for parse of one page crawling
type Instructions struct {
	ItemSelector               string `json:"itemSelector,omitempty"`
	NameOfItemSelector         string `json:"nameOfItemSelector,omitempty"`
	LinkOfItemSelector         string `json:"linkOfItemSelector,omitempty"`
	PreviewImageOfItemSelector string `json:"previewImageOfSelector,omitempty"`
	PriceOfItemSelector        string `json:"priceOfItemSelector,omitempty"`
}

type Request struct {
	IRI          string
	Instructions Instructions
}

type Response struct{ Error, Data, Message string }

func Handle(req []byte) string {
	request := Request{}

	err := json.Unmarshal(req, &request)
	if err != nil {
		warning := fmt.Sprintf(
			"Unmarshal request error: %v. Error: %v", request, err)
		fmt.Println(warning)
	}

	productsFromPage, err := pageParse(request.IRI, request.Instructions)
	if err != nil {
		warning := fmt.Sprintf(
			"Parse page error by IRI: %v. Error: %v",
			request.IRI,
			err)

		fmt.Println(warning)

		encodedResponse := Response{
			Message: warning,
			Data:    string(req),
			Error:   err.Error()}

		response, err := json.Marshal(encodedResponse)
		if err != nil {
			fmt.Println(err.Error())
		}

		return string(response)
	}

	encodedProductsFromPage, err := json.Marshal(productsFromPage)
	if err != nil {
		encodedResponse := Response{
			Data:  string(req),
			Error: err.Error()}

		response, err := json.Marshal(encodedResponse)
		if err != nil {
			fmt.Println(err.Error())
		}

		return string(response)
	}

	fmt.Println(string(encodedProductsFromPage))

	encodedResponse := Response{Data: string(encodedProductsFromPage)}

	response, err := json.Marshal(encodedResponse)
	if err != nil {
		fmt.Println(err.Error())
	}

	return string(response)
}

type Price struct {
	Value    float64
	DateTime time.Time
}

type Product struct {
	Name             string
	IRI              string
	PreviewImageLink string
	Price            Price
}

func pageParse(pageIRI string, instructions Instructions) ([]Product, error) {
	collector := colly.NewCollector(colly.Async(true))

	var productsFromPage []Product

	collector.OnHTML(instructions.ItemSelector,
		func(element *colly.HTMLElement) {

			productName := element.ChildText(instructions.NameOfItemSelector)
			productIRI := element.ChildAttr(instructions.LinkOfItemSelector, "href")
			previewImageLink := element.ChildAttr(
				instructions.PreviewImageOfItemSelector, "data-original")

			product := Product{
				Name:             productName,
				IRI:              productIRI,
				PreviewImageLink: previewImageLink}

			priceOfItemValue := element.ChildText(instructions.PriceOfItemSelector)

			patternForCutPrice, err := regexp.Compile("[ ¤]*")
			if err != nil {
				warning := fmt.Sprintf(
					"Error compile pattern for cut price for URL: %v. Error: %v",
					pageIRI,
					err)

				fmt.Println(warning)
			}

			priceIsMatched := patternForCutPrice.MatchString(priceOfItemValue)

			var priceOfItem string
			if priceIsMatched {
				priceOfItem = patternForCutPrice.ReplaceAllString(priceOfItemValue, "")
			}

			priceOfItem = strings.Replace(priceOfItem, " ", "", -1)

			priceValue, err := strconv.ParseFloat(priceOfItem, 64)
			if err != nil {
				warning := fmt.Sprintf(
					"Error get price of product: %v, by IRI: %v",
					product,
					element.Request.URL)

				fmt.Println(warning)
			}

			price := Price{
				Value:    priceValue,
				DateTime: time.Now().UTC()}

			product.Price = price

			info := fmt.Sprintf("Get product: %v by iri: %v", product, element.Request.URL)
			fmt.Println(info)

			productsFromPage = append(productsFromPage, product)
		})

	collector.OnError(func(response *colly.Response, err error) {
		warning := fmt.Sprintf(
			"Request URL: %v failed with response: %v. Error: %v",
			response.Request.URL,
			response,
			err)

		fmt.Println(warning)
	})

	err := collector.Visit(pageIRI)
	if err != nil {
		warning := fmt.Sprintf(
			"Error visit URL: %v. Error: %v",
			pageIRI,
			err)

		fmt.Println(warning)

		return nil, err
	}

	collector.Wait()

	return productsFromPage, nil
}
