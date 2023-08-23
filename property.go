package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func (c *Client) ListProperties() error {
	var r PropertiesListResponse
	currentPage := 1
	for {
		response, err := c.get(requestData{EntityEndpoint: PROPERTIES_ENDPOINT,QueryParams: fmt.Sprintf("page=%d&limit=50",currentPage)})
		if err != nil {
			return err
		}
		err = json.Unmarshal([]byte(response), &r)
		if err != nil{
			return err
		}
		for _, p := range r.Content{
			fmt.Println(p.Title)
		}
		if r.Pagination.NextPage == "null"{
			break
		}
		parsedURL, err := url.Parse(r.Pagination.NextPage)
		if err != nil {
			return err
		}
		
		nextPage, _ := strconv.Atoi(parsedURL.Query().Get("page"))
		if currentPage == nextPage{
			break
		}
		currentPage = nextPage
	}
	return nil
}

