package entities

import (
	"time"

	"github.com/google/uuid"
)

// Subscription represents a subscription entity
type Subscription struct {
	ID                 uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID     uuid.UUID  `json:"organization_id" gorm:"type:uuid;not null"`
	PlanID             uuid.UUID  `json:"plan_id" gorm:"type:uuid;not null"`
	Status             string     `json:"status" gorm:"type:varchar(50);not null;default:'active'"`
	CurrentPeriodStart time.Time  `json:"current_period_start" gorm:"not null"`
	CurrentPeriodEnd   time.Time  `json:"current_period_end" gorm:"not null"`
	CancelAtPeriodEnd  bool       `json:"cancel_at_period_end" gorm:"default:false"`
	CanceledAt         *time.Time `json:"canceled_at"`
	EndedAt            *time.Time `json:"ended_at"`
	TrialStart         *time.Time `json:"trial_start"`
	TrialEnd           *time.Time `json:"trial_end"`

	// Billing information
	BillingCycle string  `json:"billing_cycle" gorm:"type:varchar(20);not null;default:'monthly'"`
	Amount       float64 `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency     string  `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`

	// Usage tracking
	MaxUsers         int     `json:"max_users" gorm:"not null"`
	CurrentUsers     int     `json:"current_users" gorm:"default:0"`
	MaxStorageGB     int     `json:"max_storage_gb" gorm:"default:0"`
	CurrentStorageGB float64 `json:"current_storage_gb" gorm:"default:0"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Plan     *Plan     `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	Invoices []Invoice `json:"invoices,omitempty" gorm:"foreignKey:SubscriptionID"`
}

// Plan represents a subscription plan
type Plan struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Description  string    `json:"description" gorm:"type:text"`
	Price        float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Currency     string    `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	BillingCycle string    `json:"billing_cycle" gorm:"type:varchar(20);not null;default:'monthly'"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	IsPopular    bool      `json:"is_popular" gorm:"default:false"`

	// Features
	MaxUsers     int `json:"max_users" gorm:"not null"`
	MaxStorageGB int `json:"max_storage_gb" gorm:"default:0"`
	MaxProjects  int `json:"max_projects" gorm:"default:0"`
	MaxAPICalls  int `json:"max_api_calls" gorm:"default:0"`

	// Features flags
	HasAdvancedReporting  bool `json:"has_advanced_reporting" gorm:"default:false"`
	HasAPIAccess          bool `json:"has_api_access" gorm:"default:false"`
	HasCustomIntegrations bool `json:"has_custom_integrations" gorm:"default:false"`
	HasPrioritySupport    bool `json:"has_priority_support" gorm:"default:false"`
	HasPhoneSupport       bool `json:"has_phone_support" gorm:"default:false"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// Invoice represents a billing invoice
type Invoice struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;not null"`
	SubscriptionID uuid.UUID `json:"subscription_id" gorm:"type:uuid;not null"`
	InvoiceNumber  string    `json:"invoice_number" gorm:"type:varchar(50);not null;unique"`
	Status         string    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Amount         float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Currency       string    `json:"currency" gorm:"type:varchar(3);not null;default:'USD'"`
	TaxAmount      float64   `json:"tax_amount" gorm:"type:decimal(10,2);default:0"`
	DiscountAmount float64   `json:"discount_amount" gorm:"type:decimal(10,2);default:0"`
	TotalAmount    float64   `json:"total_amount" gorm:"type:decimal(10,2);not null"`

	// Billing period
	PeriodStart time.Time  `json:"period_start" gorm:"not null"`
	PeriodEnd   time.Time  `json:"period_end" gorm:"not null"`
	DueDate     time.Time  `json:"due_date" gorm:"not null"`
	PaidAt      *time.Time `json:"paid_at"`

	// Payment information
	PaymentMethodID *uuid.UUID `json:"payment_method_id" gorm:"type:uuid"`
	PaymentIntentID string     `json:"payment_intent_id" gorm:"type:varchar(100)"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	Subscription *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;not null"`
	Type           string    `json:"type" gorm:"type:varchar(20);not null"` // card, paypal, etc.
	Brand          string    `json:"brand" gorm:"type:varchar(20)"`         // visa, mastercard, etc.
	Last4          string    `json:"last4" gorm:"type:varchar(4)"`
	ExpiryMonth    int       `json:"expiry_month" gorm:"type:int"`
	ExpiryYear     int       `json:"expiry_year" gorm:"type:int"`
	Email          string    `json:"email" gorm:"type:varchar(255)"` // for paypal
	IsDefault      bool      `json:"is_default" gorm:"default:false"`
	IsActive       bool      `json:"is_active" gorm:"default:true"`

	// External payment provider ID
	ProviderID string `json:"provider_id" gorm:"type:varchar(100)"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}

// Usage represents usage tracking
type Usage struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	OrganizationID uuid.UUID `json:"organization_id" gorm:"type:uuid;not null"`
	SubscriptionID uuid.UUID `json:"subscription_id" gorm:"type:uuid;not null"`
	Metric         string    `json:"metric" gorm:"type:varchar(50);not null"` // users, storage, api_calls, etc.
	Value          float64   `json:"value" gorm:"type:decimal(10,2);not null"`
	Limit          float64   `json:"limit" gorm:"type:decimal(10,2);not null"`
	Period         string    `json:"period" gorm:"type:varchar(20);not null"` // daily, monthly, etc.
	Date           time.Time `json:"date" gorm:"not null"`

	// Audit fields
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
}
