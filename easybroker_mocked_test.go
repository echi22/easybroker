package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createMockServer(handler http.HandlerFunc) (*http.Client, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := &http.Client{
		Transport: http.DefaultTransport,
	}
	return client, server
}

func TestClient_makeRequest_Success(t *testing.T) {
	os.Setenv("API_KEY","1")
	client, server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"response": "success"}`))
	})
	defer server.Close()

	easyBrokerClient := NewClient(server.URL)
	easyBrokerClient.httpClient = client

	// Test makeRequest for success
	data := requestData{EntityEndpoint: "example"}
	response, err := easyBrokerClient.makeRequest(data, http.MethodGet)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedResponse := `{"response": "success"}`
	if response != expectedResponse {
		t.Errorf("Expected response: %s, but got: %s", expectedResponse, response)
	}
}

func TestClient_makeRequest_Failure(t *testing.T) {
	client, server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	})
	defer server.Close()

	easyBrokerClient := NewClient(server.URL)
	easyBrokerClient.httpClient = client

	// Test makeRequest for failure
	data := requestData{EntityEndpoint: "example"}
	_, err := easyBrokerClient.makeRequest(data, http.MethodGet)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestClient_ListProperties_Success_Simple(t *testing.T) {
	os.Setenv("API_KEY","1")

	client, server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pagination": {"limit": 10, "page": 1, "total": 1, "next_page": "http://test.com?page=1"}, "content": [{"agent": "Agent", "public_id": "123", "title": "Property 1"}]}`))
	})
	defer server.Close()

	easyBrokerClient := NewClient(server.URL)
	easyBrokerClient.httpClient = client

	// Test ListProperties for success
	err := easyBrokerClient.ListProperties()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestClient_ListProperties_Success_Full(t *testing.T) {
	os.Setenv("API_KEY","1")

	client, server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pagination": {"limit": 20, "page": 1, "total": 1, "next_page": "http://test.com?page=1"},
			"content": [
				{"agent": "Agent 1", "public_id": "123", "title": "Property 1", "title_image_full": "image1.jpg", "title_image_thumb": "thumb1.jpg", "bedrooms": 3, "bathrooms": 2, "parking_spaces": 1, "location": "Location 1", "property_type": "Apartment", "updated_at": "2023-08-23T12:00:00Z", "show_prices": true, "share_commission": false,
				 "operations": [
					{"type": "sale", "amount": 500000, "formated_amount": "US$ 500,000", "currency": "USD", "unit": "total", "commission": {"type": "amount", "value": 10000, "currency": "USD"}},
					{"type": "temporary_rental", "amount": 1500, "formated_amount": "US$ 1,500", "currency": "USD", "period": "monthly"}
				]}
			]
		}`))
	})
	defer server.Close()

	easyBrokerClient := NewClient(server.URL)
	easyBrokerClient.httpClient = client

	// Test ListProperties for success
	err := easyBrokerClient.ListProperties()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestClient_ListProperties_Pagination_Success(t *testing.T) {
	os.Setenv("API_KEY","1")

	client, server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		nextPage := 0
		switch page {
		case "1":
			nextPage = 2
		case "2":
			nextPage = 3
		default:
			nextPage = 0
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"pagination": {"limit": 10, "page": %s, "total": 30, "next_page": "http://test.com?page=%d"},
			"content": [
				{"agent": "Agent %s", "public_id": "%s", "title": "Property %s", "title_image_full": "image%s.jpg", "title_image_thumb": "thumb%s.jpg", "bedrooms": %d, "bathrooms": %d, "parking_spaces": %d, "location": "Location %s", "property_type": "Type %s", "updated_at": "2023-08-23T12:00:00Z", "show_prices": true, "share_commission": false,
				 "operations": [
					{"type": "sale", "amount": 500000, "formated_amount": "US$ 500,000", "currency": "USD", "unit": "total", "commission": {"type": "amount", "value": 10000, "currency": "USD"}},
					{"type": "temporary_rental", "amount": 1500, "formated_amount": "US$ 1,500", "currency": "USD", "period": "monthly"}
				]}
			]
		}`, page, nextPage, page, page, page, page, page, 3, 2, 1, page, page)))
	})
	defer server.Close()

	easyBrokerClient := NewClient(server.URL)
	easyBrokerClient.httpClient = client

	// Test ListProperties for pagination success
	err := easyBrokerClient.ListProperties()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}


func TestClient_ListProperties_Error(t *testing.T) {
	client, server := createMockServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	})
	defer server.Close()

	easyBrokerClient := NewClient(server.URL)
	easyBrokerClient.httpClient = client

	// Test ListProperties for error
	err := easyBrokerClient.ListProperties()
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}
}

func TestClient_getBaseUrl(t *testing.T) {
	// Test getBaseUrl
	easyBrokerClient := NewClient("http://example.com")
	baseURL := easyBrokerClient.getBaseUrl()
	if baseURL != "http://example.com" {
		t.Errorf("Expected base URL 'http://example.com', but got '%s'", baseURL)
	}
}

func TestClient_getHttpClient(t *testing.T) {
	// Test getHttpClient
	easyBrokerClient := NewClient("http://example.com")
	httpClient := easyBrokerClient.getHttpClient()
	if httpClient == nil {
		t.Errorf("Expected non-nil httpClient, but got nil")
	}
	// You can further test the properties of the httpClient if needed
}


func TestClient_bodyToJson(t *testing.T) {
	// Test bodyToJson
	easyBrokerClient := NewClient("http://example.com")
	bodyData := struct {
		Field string `json:"field"`
	}{Field: "value"}

	bodyReader := easyBrokerClient.bodyToJson(bodyData)
	if bodyReader == nil {
		t.Errorf("Expected non-nil bodyReader, but got nil")
	}
	// Add more tests for different scenarios
}
