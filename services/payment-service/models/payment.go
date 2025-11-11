package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type Payment struct {
	ID                    int64      `db:"id" json:"id" swaggerignore:"true"`
	OrderID               int64      `db:"order_id" json:"orderId" example:"1"`
	UserID                int64      `db:"user_id" json:"userId" swaggerignore:"true"`
	AmountCents           int        `db:"amount_cents" json:"amountCents" example:"5999"`
	Currency              string     `db:"currency" json:"currency" example:"EUR"`
	Status                string     `db:"status" json:"status" example:"pending"`
	StripePaymentIntentID *string    `db:"stripe_payment_intent_id" json:"stripePaymentIntentId,omitempty" example:"pi_1234567890"`
	StripeClientSecret    *string    `db:"stripe_client_secret" json:"stripeClientSecret,omitempty"`
	CreatedAt             time.Time  `db:"created_at" json:"createdAt" swaggerignore:"true"`
	UpdatedAt             *time.Time `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`
}

// Create creates a new payment record in the database
// used in: handlers.CreatePaymentIntent
func Create(payment *Payment) error {
	query := `INSERT INTO payments (order_id, user_id, amount_cents, currency, status, stripe_payment_intent_id, stripe_client_secret, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, now())
	          RETURNING id, created_at`

	err := db.DB.QueryRow(db.Ctx, query,
		payment.OrderID,
		payment.UserID,
		payment.AmountCents,
		payment.Currency,
		payment.Status,
		payment.StripePaymentIntentID,
		payment.StripeClientSecret,
	).Scan(&payment.ID, &payment.CreatedAt)

	return err
}

// GetByID retrieves a payment by its ID
// used in: handlers.GetPaymentStatus
func GetByID(paymentID int64) (*Payment, error) {
	payment := &Payment{}
	query := `SELECT id, order_id, user_id, amount_cents, currency, status, 
	          stripe_payment_intent_id, stripe_client_secret, created_at, updated_at
	          FROM payments WHERE id = $1`

	err := db.DB.QueryRow(db.Ctx, query, paymentID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.AmountCents,
		&payment.Currency,
		&payment.Status,
		&payment.StripePaymentIntentID,
		&payment.StripeClientSecret,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetByOrderID retrieves a payment by order ID
// used in: handlers.CreatePaymentIntent, order-service integration
func GetByOrderID(orderID int64) (*Payment, error) {
	payment := &Payment{}
	query := `SELECT id, order_id, user_id, amount_cents, currency, status, 
	          stripe_payment_intent_id, stripe_client_secret, created_at, updated_at
	          FROM payments WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1`

	err := db.DB.QueryRow(db.Ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.AmountCents,
		&payment.Currency,
		&payment.Status,
		&payment.StripePaymentIntentID,
		&payment.StripeClientSecret,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

// GetByStripePaymentIntentID retrieves a payment by Stripe Payment Intent ID
// used in: handlers.WebhookHandler
func GetByStripePaymentIntentID(paymentIntentID string) (*Payment, error) {
	payment := &Payment{}
	query := `SELECT id, order_id, user_id, amount_cents, currency, status, 
	          stripe_payment_intent_id, stripe_client_secret, created_at, updated_at
	          FROM payments WHERE stripe_payment_intent_id = $1`

	err := db.DB.QueryRow(db.Ctx, query, paymentIntentID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.UserID,
		&payment.AmountCents,
		&payment.Currency,
		&payment.Status,
		&payment.StripePaymentIntentID,
		&payment.StripeClientSecret,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return payment, nil
}

// UpdateStatus updates the payment status
// used in: handlers.ConfirmPayment, handlers.WebhookHandler
func UpdateStatus(paymentID int64, status string) error {
	query := `UPDATE payments SET status = $1, updated_at = now() WHERE id = $2`
	_, err := db.DB.Exec(db.Ctx, query, status, paymentID)
	return err
}

// GetAllByUserID retrieves all payments for a specific user
// used in: handlers.GetUserPayments (optional, für Zahlungshistorie)
func GetAllByUserID(userID int64) ([]Payment, error) {
	query := `SELECT id, order_id, user_id, amount_cents, currency, status, 
	          stripe_payment_intent_id, stripe_client_secret, created_at, updated_at
	          FROM payments WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := db.DB.Query(db.Ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment
	for rows.Next() {
		var payment Payment
		err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.UserID,
			&payment.AmountCents,
			&payment.Currency,
			&payment.Status,
			&payment.StripePaymentIntentID,
			&payment.StripeClientSecret,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, rows.Err()
}
