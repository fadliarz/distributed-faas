package domain

import (
	"testing"
)

func TestChargeDomainService_ValidateCharge(t *testing.T) {
	service := NewChargeDomainService()

	// Test valid charge
	t.Run("valid charge", func(t *testing.T) {
		charge := NewCharge(
			NewUserID("user123"),
			NewServiceID("service456"),
			NewAmount(100),
		)

		err := service.ValidateCharge(charge)
		if err != nil {
			t.Errorf("Expected no error for valid charge, got %v", err)
		}
	})

	// Test invalid charge - nil charge
	t.Run("nil charge", func(t *testing.T) {
		err := service.ValidateCharge(nil)
		if err != ErrInvalidChargeAmount {
			t.Errorf("Expected ErrInvalidChargeAmount for nil charge, got %v", err)
		}
	})

	// Test invalid charge - zero amount
	t.Run("zero amount", func(t *testing.T) {
		charge := NewCharge(
			NewUserID("user123"),
			NewServiceID("service456"),
			NewAmount(0),
		)

		err := service.ValidateCharge(charge)
		if err != ErrInvalidChargeAmount {
			t.Errorf("Expected ErrInvalidChargeAmount for zero amount, got %v", err)
		}
	})

	// Test invalid charge - negative amount
	t.Run("negative amount", func(t *testing.T) {
		charge := NewCharge(
			NewUserID("user123"),
			NewServiceID("service456"),
			NewAmount(-50),
		)

		err := service.ValidateCharge(charge)
		if err != ErrInvalidChargeAmount {
			t.Errorf("Expected ErrInvalidChargeAmount for negative amount, got %v", err)
		}
	})

	// Test invalid charge - empty user ID
	t.Run("empty user ID", func(t *testing.T) {
		charge := NewCharge(
			NewUserID(""),
			NewServiceID("service456"),
			NewAmount(100),
		)

		err := service.ValidateCharge(charge)
		if err != ErrInvalidUserID {
			t.Errorf("Expected ErrInvalidUserID for empty user ID, got %v", err)
		}
	})

	// Test invalid charge - empty service ID
	t.Run("empty service ID", func(t *testing.T) {
		charge := NewCharge(
			NewUserID("user123"),
			NewServiceID(""),
			NewAmount(100),
		)

		err := service.ValidateCharge(charge)
		if err != ErrInvalidServiceID {
			t.Errorf("Expected ErrInvalidServiceID for empty service ID, got %v", err)
		}
	})
}
