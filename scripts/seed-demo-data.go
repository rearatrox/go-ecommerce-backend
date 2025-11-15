package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	baseURL = "http://localhost"
)

var (
	userServicePort    = getEnv("USERSERVICE_PORT", "8081")
	productServicePort = getEnv("PRODUCTSERVICE_PORT", "8082")
	cartServicePort    = getEnv("CARTSERVICE_PORT", "8083")
	apiPrefix          = getEnv("API_PREFIX", "/api/v1")
)

// User credentials
var (
	adminToken     string
	customer1Token string
	customer2Token string
)

func main() {
	fmt.Println("🌱 Starting Demo Data Seeding...")
	fmt.Println("=====================================")

	// Wait for services to be ready
	waitForServices()

	// 1. Create Demo Users
	fmt.Println("\n📝 Creating demo users...")
	createDemoUsers()

	// 2. Login and get tokens
	fmt.Println("\n🔐 Logging in users...")
	loginUsers()

	// 3. Create Addresses
	fmt.Println("\n🏠 Creating addresses...")
	createAddresses()

	// 4. Create Categories
	fmt.Println("\n📁 Creating categories...")
	createCategories()

	// 5. Create Products
	fmt.Println("\n📦 Creating products...")
	createProducts()

	// 6. Add items to carts
	fmt.Println("\n🛒 Adding items to carts...")
	addItemsToCarts()

	fmt.Println("\n=====================================")
	fmt.Println("✅ Demo data seeding completed!")
	fmt.Println("\n📊 Demo Users Created:")
	fmt.Println("  Admin:")
	fmt.Println("    Email: Admin@example.com")
	fmt.Println("    Password: Admin123!")
	fmt.Println("  Customer 1:")
	fmt.Println("    Email: customer1@example.com")
	fmt.Println("    Password: Customer123!")
	fmt.Println("  Customer 2:")
	fmt.Println("    Email: customer2@example.com")
	fmt.Println("    Password: Customer123!")
}

func waitForServices() {
	fmt.Println("⏳ Waiting for services to be ready...")

	// Try to login admin to check if user-service is ready
	// This also verifies that migrations have run (admin user exists)
	fmt.Println("  Checking user-service...")
	for i := 0; i < 30; i++ {
		loginData := map[string]string{
			"email":    "Admin@example.com",
			"password": "Admin123!",
		}
		jsonData, _ := json.Marshal(loginData)
		url := fmt.Sprintf("%s:%s%s/auth/login", baseURL, userServicePort, apiPrefix)

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err == nil && (resp.StatusCode == 200 || resp.StatusCode == 401 || resp.StatusCode == 404) {
			// 200 = login success, 401/404 = service is up (just wrong credentials or user doesn't exist yet)
			resp.Body.Close()
			fmt.Printf("  ✓ user-service ready\n")
			break
		}
		if err == nil {
			resp.Body.Close()
		}
		if i == 29 {
			fmt.Printf("  ⚠ user-service might not be ready, continuing anyway...\n")
		}
		time.Sleep(1 * time.Second)
	}

	// Just wait a moment for other services - they'll error later if not ready
	fmt.Println("  ✓ Assuming other services are ready")
	time.Sleep(2 * time.Second)
}

func createDemoUsers() {
	users := []map[string]interface{}{
		{
			"email":     "customer1@example.com",
			"password":  "Customer123!",
			"firstName": "John",
			"lastName":  "Doe",
			"phone":     "+49123456789",
		},
		{
			"email":     "customer2@example.com",
			"password":  "Customer123!",
			"firstName": "Jane",
			"lastName":  "Smith",
			"phone":     "+49987654321",
		},
	}

	url := fmt.Sprintf("%s:%s%s/auth/signup", baseURL, userServicePort, apiPrefix)

	for _, user := range users {
		jsonData, _ := json.Marshal(user)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("  ✗ Failed to create user %s: %v\n", user["email"], err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 201 || resp.StatusCode == 409 { // 409 = already exists
			fmt.Printf("  ✓ User created/exists: %s\n", user["email"])
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("  ✗ Failed to create user %s: %s\n", user["email"], string(body))
		}
	}
}

func loginUsers() {
	credentials := []struct {
		email    string
		password string
		token    *string
		role     string
	}{
		{"admin@example.com", "admin123", &adminToken, "admin"},
		{"customer1@example.com", "Customer123!", &customer1Token, "customer1"},
		{"customer2@example.com", "Customer123!", &customer2Token, "customer2"},
	}

	url := fmt.Sprintf("%s:%s%s/auth/login", baseURL, userServicePort, apiPrefix)

	for _, cred := range credentials {
		loginData := map[string]string{
			"email":    cred.email,
			"password": cred.password,
		}
		jsonData, _ := json.Marshal(loginData)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("  ✗ Failed to login %s: %v\n", cred.email, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			var result map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&result)
			if token, ok := result["token"].(string); ok {
				*cred.token = token
				fmt.Printf("  ✓ Logged in: %s\n", cred.email)
			}
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("  ✗ Failed to login %s: %s\n", cred.email, string(body))
		}
	}
}

func createAddresses() {
	addresses := []struct {
		token   string
		user    string
		address map[string]interface{}
	}{
		{
			customer1Token,
			"customer1",
			map[string]interface{}{
				"fullName":   "Max Mustermann",
				"street":     "Musterstraße 123",
				"city":       "Berlin",
				"postalCode": "10115",
				"country":    "Germany",
				"isDefault":  true,
				"type":       "shipping",
			},
		},
		{
			customer1Token,
			"customer1",
			map[string]interface{}{
				"fullName":   "Max Mustermann",
				"street":     "Beispielweg 45",
				"city":       "München",
				"postalCode": "80331",
				"country":    "Germany",
				"isDefault":  false,
				"type":       "billing",
			},
		},
		{
			customer2Token,
			"customer2",
			map[string]interface{}{
				"fullName":   "Max Mustermann",
				"street":     "Teststraße 789",
				"city":       "Hamburg",
				"postalCode": "20095",
				"country":    "Germany",
				"isDefault":  true,
				"type":       "shipping",
			},
		},
		{
			customer2Token,
			"customer2",
			map[string]interface{}{
				"fullName":   "Max Mustermann",
				"street":     "Teststraße 789",
				"city":       "Hamburg",
				"postalCode": "20095",
				"country":    "Germany",
				"isDefault":  true,
				"type":       "billing",
			},
		},
	}

	url := fmt.Sprintf("%s:%s%s/users/me/addresses", baseURL, userServicePort, apiPrefix)

	for _, addr := range addresses {
		jsonData, _ := json.Marshal(addr.address)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", addr.token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("  ✗ Failed to create address for %s: %v\n", addr.user, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 201 {
			fmt.Printf("  ✓ Address created for %s: %s, %s\n", addr.user, addr.address["city"], addr.address["type"])
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("  ✗ Failed to create address for %s: %s\n", addr.user, string(body))
		}
	}
}

func createCategories() {
	categories := []map[string]interface{}{
		{
			"name":        "Electronics",
			"slug":        "electronics",
			"description": "Electronic devices and gadgets",
		},
		{
			"name":        "Clothing",
			"slug":        "clothing",
			"description": "Fashion and apparel",
		},
		{
			"name":        "Books",
			"slug":        "books",
			"description": "Books and literature",
		},
		{
			"name":        "Home & Garden",
			"slug":        "home-garden",
			"description": "Home decoration and garden supplies",
		},
	}

	url := fmt.Sprintf("%s:%s%s/admin/categories/create", baseURL, productServicePort, apiPrefix)

	for _, category := range categories {
		jsonData, _ := json.Marshal(category)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", adminToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("  ✗ Failed to create category %s: %v\n", category["name"], err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 201 || resp.StatusCode == 409 {
			fmt.Printf("  ✓ Category created/exists: %s\n", category["name"])
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("  ✗ Failed to create category %s: %s\n", category["name"], string(body))
		}
	}
}

func createProducts() {
	products := []map[string]interface{}{
		{
			"name":        "Gaming Laptop XPS 15",
			"sku":         "LAPTOP-001",
			"description": "High-performance gaming laptop with RTX 4070",
			"priceCents":  149999,
			"stockQty":    25,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1603302576837-37561b2e2302?w=800",
			"currency":    "EUR",
			"categoryIds": []int{1}, // Electronics
		},
		{
			"name":        "Wireless Mouse MX Master 3",
			"sku":         "MOUSE-001",
			"description": "Ergonomic wireless mouse with precision tracking",
			"priceCents":  9999,
			"stockQty":    100,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=800",
			"currency":    "EUR",
			"categoryIds": []int{1},
		},
		{
			"name":        "USB-C Charging Cable",
			"sku":         "CABLE-001",
			"description": "2-meter braided USB-C charging cable",
			"priceCents":  1990,
			"stockQty":    250,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1591290619762-d71b02ae3c99?w=800",
			"currency":    "EUR",
			"categoryIds": []int{1},
		},
		{
			"name":        "Cotton T-Shirt Blue",
			"sku":         "TSHIRT-001",
			"description": "Comfortable 100% cotton t-shirt in classic blue",
			"priceCents":  2999,
			"stockQty":    150,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=800",
			"currency":    "EUR",
			"categoryIds": []int{2}, // Clothing
		},
		{
			"name":        "Slim Fit Jeans",
			"sku":         "JEANS-001",
			"description": "Modern slim fit jeans in dark wash",
			"priceCents":  7999,
			"stockQty":    75,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1542272604-787c3835535d?w=800",
			"currency":    "EUR",
			"categoryIds": []int{2},
		},
		{
			"name":        "The Go Programming Language",
			"sku":         "BOOK-001",
			"description": "Comprehensive guide to Go programming by experts",
			"priceCents":  4999,
			"stockQty":    50,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1532012197267-da84d127e765?w=800",
			"currency":    "EUR",
			"categoryIds": []int{3}, // Books
		},
		{
			"name":        "Clean Code",
			"sku":         "BOOK-002",
			"description": "A handbook of agile software craftsmanship by Robert C. Martin",
			"priceCents":  4499,
			"stockQty":    40,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1544947950-fa07a98d237f?w=800",
			"currency":    "EUR",
			"categoryIds": []int{3},
		},
		{
			"name":        "LED Desk Lamp",
			"sku":         "LAMP-001",
			"description": "Adjustable LED desk lamp with touch controls",
			"priceCents":  5999,
			"stockQty":    60,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1507473885765-e6ed057f782c?w=800",
			"currency":    "EUR",
			"categoryIds": []int{4}, // Home & Garden
		},
		{
			"name":        "Ceramic Plant Pot",
			"sku":         "POT-001",
			"description": "Modern 20cm ceramic plant pot with drainage",
			"priceCents":  1999,
			"stockQty":    120,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1485955900006-10f4d324d411?w=800",
			"currency":    "EUR",
			"categoryIds": []int{4},
		},
		{
			"name":        "Coffee Mug Set (4pc)",
			"sku":         "MUG-001",
			"description": "Set of 4 premium ceramic coffee mugs",
			"priceCents":  2999,
			"stockQty":    80,
			"status":      "active",
			"imageUrl":    "https://images.unsplash.com/photo-1514228742587-6b1558fcca3d?w=800",
			"currency":    "EUR",
			"categoryIds": []int{4},
		},
	}

	url := fmt.Sprintf("%s:%s%s/admin/products/create", baseURL, productServicePort, apiPrefix)

	for _, product := range products {
		jsonData, _ := json.Marshal(product)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", adminToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("  ✗ Failed to create product %s: %v\n", product["name"], err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 201 || resp.StatusCode == 409 {
			fmt.Printf("  ✓ Product created/exists: %s (€%.2f)\n", product["name"], float64(product["priceCents"].(int))/100)
		} else {
			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("  ✗ Failed to create product %s: %s\n", product["name"], string(body))
		}
	}
}

func addItemsToCarts() {
	cartItems := []struct {
		token string
		user  string
		items []map[string]interface{}
	}{
		{
			customer1Token,
			"customer1",
			[]map[string]interface{}{
				{"productId": 1, "quantity": 1}, // Laptop
				{"productId": 2, "quantity": 2}, // Mouse
				{"productId": 6, "quantity": 1}, // Book
			},
		},
		{
			customer2Token,
			"customer2",
			[]map[string]interface{}{
				{"productId": 4, "quantity": 2},  // T-Shirt
				{"productId": 10, "quantity": 1}, // Coffee Mug Set
			},
		},
	}

	url := fmt.Sprintf("%s:%s%s/cart/items", baseURL, cartServicePort, apiPrefix)

	for _, cart := range cartItems {
		for _, item := range cart.items {
			jsonData, _ := json.Marshal(item)
			req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", cart.token)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("  ✗ Failed to add item to cart for %s: %v\n", cart.user, err)
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == 201 || resp.StatusCode == 200 {
				fmt.Printf("  ✓ Item added to cart for %s: Product ID %v (x%v)\n",
					cart.user, item["productId"], item["quantity"])
			} else {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("  ✗ Failed to add item to cart for %s: %s\n", cart.user, string(body))
			}
		}
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
