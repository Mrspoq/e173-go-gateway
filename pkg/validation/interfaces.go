package validation

import "github.com/e173-gateway/e173_go_gateway/pkg/models"

// PhoneNumberValidator interface for phone number validation
type PhoneNumberValidator interface {
	IsValid(phoneNumber string) bool
}

// WhatsAppValidator interface for WhatsApp validation
type WhatsAppValidator interface {
	ValidateNumber(phoneNumber string) (*models.ValidationResult, error)
}

// Ensure implementations
var _ PhoneNumberValidator = (*GooglePhoneValidator)(nil)
var _ PhoneNumberValidator = (*LibPhoneNumberValidator)(nil)
var _ WhatsAppValidator = (*WhatsAppBusinessValidator)(nil)