package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type OrderResponse struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"userId"`
	Status     string `json:"status"`
	TotalCents int    `json:"totalCents"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

// getOrderDetails fetches order information from order-service
func getOrderDetails(orderID int64) (*OrderResponse, error) {
	orderServiceURL := "http://order-service:8080"

	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	url := fmt.Sprintf("%s%s/orders/%d", orderServiceURL, apiPrefix, orderID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call order-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("order-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var order OrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	return &order, nil
}

// updateOrderStatus updates the order status in order-service
func updateOrderStatus(orderID int64, status string) error {
	orderServiceURL := "http://order-service:8080"

	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	url := fmt.Sprintf("%s%s/internal/orders/%d/status", orderServiceURL, apiPrefix, orderID)

	reqBody := UpdateOrderStatusRequest{
		Status: status,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call order-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("order-service returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
