#!/bin/bash

echo "Generating Swagger documentation for all services..."

services=("user-service" "product-service" "cart-service" "order-service")

for service in "${services[@]}"; do
    echo "Generating Swagger docs for $service..."
    cd "services/$service"
    swag init
    cd "../.."
    echo "✅ $service done"
done

echo "🎉 All Swagger documentation generated!"