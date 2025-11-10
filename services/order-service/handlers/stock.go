package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"rearatrox/go-ecommerce-backend/services/order-service/models"
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

type ReduceStockRequest struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

// checkStockAvailability calls the product-service to check stock
func checkStockAvailability(productID int64, quantity int) (*CheckStockResponse, error) {
	port := os.Getenv("PRODUCTSERVICE_PORT")
	if port == "" {
		port = "8081" // Default product-service port
	}
	productServiceURL := fmt.Sprintf("http://product-service:%s", port)

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

// reduceStock calls the product-service to reduce stock (with JWT token)
func reduceStock(productID int64, quantity int, jwtToken string) error {
	port := os.Getenv("PRODUCTSERVICE_PORT")
	if port == "" {
		port = "8080" // Default internal container port
	}
	productServiceURL := fmt.Sprintf("http://product-service:%s", port)

	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	url := fmt.Sprintf("%s%s/admin/products/stock/reduce", productServiceURL, apiPrefix)

	reqBody := ReduceStockRequest{
		ProductID: productID,
		Quantity:  quantity,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stock reduce failed: %s", string(body))
	}

	return nil
}

// reduceStockForOrder reduces stock for all items in an order
func reduceStockForOrder(items []models.OrderItem, jwtToken string) error {
	for _, item := range items {
		if err := reduceStock(item.ProductID, item.Quantity, jwtToken); err != nil {
			return fmt.Errorf("failed to reduce stock for product %d: %w", item.ProductID, err)
		}
	}
	return nil
}
