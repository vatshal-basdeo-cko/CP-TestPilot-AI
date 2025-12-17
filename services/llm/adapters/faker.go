package adapters

import (
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

// FakerAdapter generates realistic test data
type FakerAdapter struct {
	faker *gofakeit.Faker
}

// NewFakerAdapter creates a new faker adapter
func NewFakerAdapter() *FakerAdapter {
	return &FakerAdapter{
		faker: gofakeit.New(0), // Random seed
	}
}

// GenerateByType generates data based on field type and format
func (f *FakerAdapter) GenerateByType(fieldName, fieldType, format string) interface{} {
	// Normalize inputs
	fieldName = strings.ToLower(fieldName)
	fieldType = strings.ToLower(fieldType)
	format = strings.ToLower(format)

	// Check format first
	switch format {
	case "email":
		return f.faker.Email()
	case "phone":
		return f.faker.Phone()
	case "date":
		return f.faker.Date().Format("2006-01-02")
	case "datetime":
		return f.faker.Date().Format("2006-01-02T15:04:05Z")
	case "uuid":
		return f.faker.UUID()
	case "url":
		return f.faker.URL()
	case "ipv4":
		return f.faker.IPv4Address()
	case "ipv6":
		return f.faker.IPv6Address()
	case "card", "credit_card", "creditcard":
		return f.generateCreditCard()
	case "cvv":
		return f.faker.CreditCardCvv()
	case "expiry", "card_expiry":
		return f.faker.CreditCardExp()
	case "currency":
		return f.faker.CurrencyShort()
	}

	// Check field name patterns
	switch {
	case strings.Contains(fieldName, "email"):
		return f.faker.Email()
	case strings.Contains(fieldName, "phone"):
		return f.faker.Phone()
	case strings.Contains(fieldName, "first_name") || strings.Contains(fieldName, "firstname"):
		return f.faker.FirstName()
	case strings.Contains(fieldName, "last_name") || strings.Contains(fieldName, "lastname"):
		return f.faker.LastName()
	case strings.Contains(fieldName, "name"):
		return f.faker.Name()
	case strings.Contains(fieldName, "address"):
		return f.faker.Address().Address
	case strings.Contains(fieldName, "city"):
		return f.faker.City()
	case strings.Contains(fieldName, "country"):
		return f.faker.Country()
	case strings.Contains(fieldName, "zip") || strings.Contains(fieldName, "postal"):
		return f.faker.Zip()
	case strings.Contains(fieldName, "card") || strings.Contains(fieldName, "pan"):
		return f.generateCreditCard()
	case strings.Contains(fieldName, "cvv") || strings.Contains(fieldName, "cvc"):
		return f.faker.CreditCardCvv()
	case strings.Contains(fieldName, "expir"):
		return f.faker.CreditCardExp()
	case strings.Contains(fieldName, "amount") || strings.Contains(fieldName, "price"):
		return f.faker.Float64Range(10.0, 1000.0)
	case strings.Contains(fieldName, "currency"):
		return f.faker.CurrencyShort()
	case strings.Contains(fieldName, "date"):
		return f.faker.Date().Format("2006-01-02")
	case strings.Contains(fieldName, "description"):
		return f.faker.Sentence(10)
	case strings.Contains(fieldName, "id"):
		return f.faker.UUID()
	case strings.Contains(fieldName, "company"):
		return f.faker.Company()
	case strings.Contains(fieldName, "url") || strings.Contains(fieldName, "website"):
		return f.faker.URL()
	}

	// Default based on type
	switch fieldType {
	case "string":
		return f.faker.LoremIpsumWord()
	case "number", "integer", "int":
		return f.faker.IntRange(1, 1000)
	case "float", "decimal", "double":
		return f.faker.Float64Range(0.01, 1000.0)
	case "boolean", "bool":
		return f.faker.Bool()
	case "array":
		return []string{f.faker.LoremIpsumWord()}
	case "object":
		return map[string]interface{}{"key": f.faker.LoremIpsumWord()}
	default:
		return f.faker.LoremIpsumWord()
	}
}

// generateCreditCard generates a test credit card number
func (f *FakerAdapter) generateCreditCard() string {
	// Use a well-known test card number pattern
	cards := []string{
		"4111111111111111", // Visa test
		"5555555555554444", // Mastercard test
		"378282246310005",  // Amex test
	}
	return cards[f.faker.IntRange(0, len(cards)-1)]
}

// GenerateMultiple generates multiple values for a field
func (f *FakerAdapter) GenerateMultiple(fieldName, fieldType, format string, count int) []interface{} {
	results := make([]interface{}, count)
	for i := 0; i < count; i++ {
		results[i] = f.GenerateByType(fieldName, fieldType, format)
	}
	return results
}

