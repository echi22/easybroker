package main

import (
	"net/http"
	"time"
)

/*
A Client is an EasyBroker Client.
BaseURL can be configured, if not, 'http://localhost:8080' will be used as default for testing porpoises.
HTTPCLient default TimeOut is set to 60 seconds.
*/
type Client struct {
	BaseURL    string
	httpClient *http.Client
}

type Commission struct {
	Type     string  `json:"type"`
	Value    interface{} `json:"value"`
	Currency string  `json:"currency"`
}

type Operation struct {
	Type           string     `json:"type"`
	Amount         float64    `json:"amount"`
	FormattedAmount string     `json:"formated_amount"`
	Currency       string     `json:"currency"`
	Unit           string     `json:"unit"`
	Commission     Commission `json:"commission"`
	Period         string     `json:"period"`
}

type Property struct {
	Agent           string      `json:"agent"`
	PublicID        string      `json:"public_id"`
	Title           string      `json:"title"`
	TitleImageFull  string      `json:"title_image_full"`
	TitleImageThumb string      `json:"title_image_thumb"`
	Bedrooms        int         `json:"bedrooms"`
	Bathrooms       int         `json:"bathrooms"`
	ParkingSpaces   int         `json:"parking_spaces"`
	Location        string      `json:"location"`
	PropertyType    string      `json:"property_type"`
	UpdatedAt       time.Time   `json:"updated_at"`
	ShowPrices      bool        `json:"show_prices"`
	ShareCommission bool        `json:"share_commission"`
	Operations      []Operation `json:"operations"`
}

type Pagination struct{
	Limit int `json:"limit"`
	Page int `json:"page"`
	Total int `json:"total"`
	NextPage string `json:"next_page"`

}
type PropertiesListResponse struct{
	Pagination Pagination `json:"pagination"`
	Content []Property `json:"content"`
}
/*
This struct is used to generate request URL and add body if needed.
*EntityEndpoint: Required. Endpoint path to add after base url. e.g http://base_url/EntityEndpoint
*Id: Optional. If request is GET, DELETE, PATCH, must be present to identify the resource. It goes after EntityEndpoint. E.g: e.g http://base_url/EntityEndpoint/Id
*QueryParams: Optional. If present, added at the end of the url. E.g: e.g http://base_url/EntityEndpoint/Id?QueryParams
*Body: Optional. If present, will be converted to JSON and sent in request body.
*/
type requestData struct {
	EntityEndpoint string      `json:"entity_endpoint"`
	Id             string      `json:"id"`
	QueryParams    string      `json:"query_params"`
	Body           interface{} `json:"body"`
}
