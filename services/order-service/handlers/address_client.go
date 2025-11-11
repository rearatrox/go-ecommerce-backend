package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type AddressResponse struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"userId"`
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
	Country string `json:"country"`
}

// verifyAddressOwnership checks if an address belongs to the given user
func verifyAddressOwnership(addressID int64, userID int64, jwtToken string) error {
	if addressID == 0 {
		return nil // No address specified, which is allowed
	}

	userServiceURL := "http://user-service:8080"

	apiPrefix := os.Getenv("API_PREFIX")
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	url := fmt.Sprintf("%s%s/users/me/addresses/%d", userServiceURL, apiPrefix, addressID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", jwtToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call user-service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("address not found")
	}

	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("address does not belong to user")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("user-service returned status %d: %s", resp.StatusCode, string(body))
	}

	var address AddressResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return fmt.Errorf("failed to decode address response: %w", err)
	}

	// Double-check that the address belongs to the user
	if address.UserID != userID {
		return fmt.Errorf("address does not belong to user")
	}

	return nil
}
