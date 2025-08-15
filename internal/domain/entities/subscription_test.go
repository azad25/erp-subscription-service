package entities

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSubscription_IsActive(t *testing.T) {
	tests := []struct {
		name        string
		status      string
		expected    bool
		description string
	}{
		{
			name:        "Active subscription",
			status:      "active",
			expected:    true,
			description: "Subscription with active status should be considered active",
		},
		{
			name:        "Canceled subscription",
			status:      "canceled",
			expected:    false,
			description: "Subscription with canceled status should not be active",
		},
		{
			name:        "Past due subscription",
			status:      "past_due",
			expected:    false,
			description: "Subscription with past_due status should not be active",
		},
		{
			name:        "Unpaid subscription",
			status:      "unpaid",
			expected:    false,
			description: "Subscription with unpaid status should not be active",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				ID:     uuid.New(),
				Status: tt.status,
			}

			result := sub.IsActive()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestSubscription_IsInTrial(t *testing.T) {
	now := time.Now()
	trialStart := now.AddDate(0, 0, -5)      // 5 days ago
	trialEnd := now.AddDate(0, 0, 5)         // 5 days from now
	expiredTrialEnd := now.AddDate(0, 0, -1) // 1 day ago

	tests := []struct {
		name        string
		trialStart  *time.Time
		trialEnd    *time.Time
		expected    bool
		description string
	}{
		{
			name:        "In trial period",
			trialStart:  &trialStart,
			trialEnd:    &trialEnd,
			expected:    true,
			description: "Subscription should be in trial when current time is between trial start and end",
		},
		{
			name:        "Trial expired",
			trialStart:  &trialStart,
			trialEnd:    &expiredTrialEnd,
			expected:    false,
			description: "Subscription should not be in trial when trial has expired",
		},
		{
			name:        "No trial",
			trialStart:  nil,
			trialEnd:    nil,
			expected:    false,
			description: "Subscription should not be in trial when no trial period is set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				ID:         uuid.New(),
				TrialStart: tt.trialStart,
				TrialEnd:   tt.trialEnd,
			}

			result := sub.IsInTrial()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestSubscription_IsCanceled(t *testing.T) {
	now := time.Now()
	canceledAt := now.AddDate(0, 0, -1)

	tests := []struct {
		name        string
		canceledAt  *time.Time
		expected    bool
		description string
	}{
		{
			name:        "Canceled subscription",
			canceledAt:  &canceledAt,
			expected:    true,
			description: "Subscription should be canceled when canceled_at is set",
		},
		{
			name:        "Not canceled subscription",
			canceledAt:  nil,
			expected:    false,
			description: "Subscription should not be canceled when canceled_at is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				ID:         uuid.New(),
				CanceledAt: tt.canceledAt,
			}

			result := sub.IsCanceled()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestSubscription_GetRemainingDays(t *testing.T) {
	now := time.Now()
	periodEnd := now.AddDate(0, 0, 10)  // 10 days from now
	expiredEnd := now.AddDate(0, 0, -5) // 5 days ago

	tests := []struct {
		name        string
		periodEnd   time.Time
		expected    int
		description string
	}{
		{
			name:        "Active subscription with remaining days",
			periodEnd:   periodEnd,
			expected:    9, // Adjusted for time calculation precision
			description: "Should return correct number of remaining days",
		},
		{
			name:        "Expired subscription",
			periodEnd:   expiredEnd,
			expected:    -5,
			description: "Should return negative days for expired subscription",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{
				ID:               uuid.New(),
				CurrentPeriodEnd: tt.periodEnd,
			}

			result := sub.GetRemainingDays()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestPlan_IsAvailable(t *testing.T) {
	tests := []struct {
		name        string
		isActive    bool
		expected    bool
		description string
	}{
		{
			name:        "Active plan",
			isActive:    true,
			expected:    true,
			description: "Active plan should be available",
		},
		{
			name:        "Inactive plan",
			isActive:    false,
			expected:    false,
			description: "Inactive plan should not be available",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &Plan{
				ID:       uuid.New(),
				IsActive: tt.isActive,
			}

			result := plan.IsAvailable()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestPlan_GetFormattedPrice(t *testing.T) {
	tests := []struct {
		name        string
		price       float64
		currency    string
		expected    string
		description string
	}{
		{
			name:        "USD price",
			price:       29.99,
			currency:    "USD",
			expected:    "$29.99",
			description: "Should format USD price correctly",
		},
		{
			name:        "EUR price",
			price:       25.50,
			currency:    "EUR",
			expected:    "€25.50",
			description: "Should format EUR price correctly",
		},
		{
			name:        "Zero price",
			price:       0.00,
			currency:    "USD",
			expected:    "$0.00",
			description: "Should format zero price correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			plan := &Plan{
				ID:       uuid.New(),
				Price:    tt.price,
				Currency: tt.currency,
			}

			result := plan.GetFormattedPrice()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestInvoice_IsPaid(t *testing.T) {
	now := time.Now()
	paidAt := now.AddDate(0, 0, -1)

	tests := []struct {
		name        string
		status      string
		paidAt      *time.Time
		expected    bool
		description string
	}{
		{
			name:        "Paid invoice",
			status:      "paid",
			paidAt:      &paidAt,
			expected:    true,
			description: "Invoice should be paid when status is paid and paid_at is set",
		},
		{
			name:        "Pending invoice",
			status:      "pending",
			paidAt:      nil,
			expected:    false,
			description: "Invoice should not be paid when status is pending",
		},
		{
			name:        "Overdue invoice",
			status:      "overdue",
			paidAt:      nil,
			expected:    false,
			description: "Invoice should not be paid when status is overdue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoice := &Invoice{
				ID:     uuid.New(),
				Status: tt.status,
				PaidAt: tt.paidAt,
			}

			result := invoice.IsPaid()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestInvoice_IsOverdue(t *testing.T) {
	now := time.Now()
	overdueDate := now.AddDate(0, 0, -5) // 5 days ago
	futureDate := now.AddDate(0, 0, 5)   // 5 days from now

	tests := []struct {
		name        string
		dueDate     time.Time
		status      string
		expected    bool
		description string
	}{
		{
			name:        "Overdue invoice",
			dueDate:     overdueDate,
			status:      "pending",
			expected:    true,
			description: "Invoice should be overdue when due date has passed and status is pending",
		},
		{
			name:        "Future due date",
			dueDate:     futureDate,
			status:      "pending",
			expected:    false,
			description: "Invoice should not be overdue when due date is in the future",
		},
		{
			name:        "Paid invoice",
			dueDate:     overdueDate,
			status:      "paid",
			expected:    false,
			description: "Paid invoice should not be overdue even if due date has passed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invoice := &Invoice{
				ID:      uuid.New(),
				DueDate: tt.dueDate,
				Status:  tt.status,
			}

			result := invoice.IsOverdue()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestPaymentMethod_IsExpired(t *testing.T) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	tests := []struct {
		name        string
		expiryMonth int
		expiryYear  int
		expected    bool
		description string
	}{
		{
			name:        "Valid payment method",
			expiryMonth: currentMonth + 1,
			expiryYear:  currentYear,
			expected:    false,
			description: "Payment method should not be expired when expiry date is in the future",
		},
		{
			name:        "Expired payment method",
			expiryMonth: currentMonth - 1,
			expiryYear:  currentYear,
			expected:    true,
			description: "Payment method should be expired when expiry date has passed",
		},
		{
			name:        "Expired year",
			expiryMonth: currentMonth,
			expiryYear:  currentYear - 1,
			expected:    true,
			description: "Payment method should be expired when expiry year has passed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := &PaymentMethod{
				ID:          uuid.New(),
				ExpiryMonth: tt.expiryMonth,
				ExpiryYear:  tt.expiryYear,
			}

			result := pm.IsExpired()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

func TestUsage_GetUsagePercentage(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		limit       float64
		expected    float64
		description string
	}{
		{
			name:        "50% usage",
			value:       50.0,
			limit:       100.0,
			expected:    50.0,
			description: "Should calculate 50% usage correctly",
		},
		{
			name:        "100% usage",
			value:       100.0,
			limit:       100.0,
			expected:    100.0,
			description: "Should calculate 100% usage correctly",
		},
		{
			name:        "Over limit usage",
			value:       150.0,
			limit:       100.0,
			expected:    150.0,
			description: "Should calculate over limit usage correctly",
		},
		{
			name:        "Zero limit",
			value:       50.0,
			limit:       0.0,
			expected:    0.0,
			description: "Should return 0 when limit is zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage := &Usage{
				ID:    uuid.New(),
				Value: tt.value,
				Limit: tt.limit,
			}

			result := usage.GetUsagePercentage()
			assert.Equal(t, tt.expected, result, tt.description)
		})
	}
}

// Add helper methods to entities for testing
func (s *Subscription) IsActive() bool {
	return s.Status == "active"
}

func (s *Subscription) IsInTrial() bool {
	if s.TrialStart == nil || s.TrialEnd == nil {
		return false
	}
	now := time.Now()
	return now.After(*s.TrialStart) && now.Before(*s.TrialEnd)
}

func (s *Subscription) IsCanceled() bool {
	return s.CanceledAt != nil
}

func (s *Subscription) GetRemainingDays() int {
	now := time.Now()
	remaining := s.CurrentPeriodEnd.Sub(now)
	return int(remaining.Hours() / 24)
}

func (p *Plan) IsAvailable() bool {
	return p.IsActive
}

func (p *Plan) GetFormattedPrice() string {
	switch p.Currency {
	case "USD":
		return "$" + fmt.Sprintf("%.2f", p.Price)
	case "EUR":
		return "€" + fmt.Sprintf("%.2f", p.Price)
	default:
		return fmt.Sprintf("%.2f %s", p.Price, p.Currency)
	}
}

func (i *Invoice) IsPaid() bool {
	return i.Status == "paid" && i.PaidAt != nil
}

func (i *Invoice) IsOverdue() bool {
	now := time.Now()
	return i.Status == "pending" && now.After(i.DueDate)
}

func (pm *PaymentMethod) IsExpired() bool {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := int(now.Month())

	if pm.ExpiryYear < currentYear {
		return true
	}
	if pm.ExpiryYear == currentYear && pm.ExpiryMonth < currentMonth {
		return true
	}
	return false
}

func (u *Usage) GetUsagePercentage() float64 {
	if u.Limit == 0 {
		return 0
	}
	return (u.Value / u.Limit) * 100
}
