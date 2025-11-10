package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CheckStockRequest struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type CheckStockResponse struct {
	Available    bool  `json:"available"`
	RequestedQty int   `json:"requestedQty"`
	AvailableQty int   `json:"availableQty"`
	ProductID    int64 `json:"productId"`
}

// checkStockAvailability calls the product-service to check stock
func checkStockAvailability(productID int64, quantity int) (*CheckStockResponse, error) {
	productServiceURL := "http://product-service:" + os.Getenv("PRODUCTSERVICE_PORT")
	if productServiceURL == "" {
		productServiceURL = "http://product-service:8080" // Default for docker-compose
	}

	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	url := fmt.Sprintf("%s%s/products/stock/check", productServiceURL, apiPrefix)

	reqBody := CheckStockRequest{
		ProductID: productID,
		Quantity:  quantity,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stock check failed: %s", string(body))
	}

	var stockResp CheckStockResponse
	if err := json.NewDecoder(resp.Body).Decode(&stockResp); err != nil {
		return nil, err
	}

	return &stockResp, nil
}
