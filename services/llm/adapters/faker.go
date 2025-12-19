package adapters

import (
	"encoding/binary"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
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
	// CKO-specific ID formats (Checkout.com style: prefix + base32 GUID)
	case strings.Contains(fieldName, "payment_id") || strings.Contains(fieldName, "pay_id"):
		return f.generateCKOId("pay")
	case strings.Contains(fieldName, "action_id") || strings.Contains(fieldName, "act_id"):
		return f.generateCKOId("act")
	case strings.Contains(fieldName, "entity_id") || strings.Contains(fieldName, "ent_id"):
		return f.generateCKOId("ent")
	case strings.Contains(fieldName, "transaction_id") || strings.Contains(fieldName, "txn_id"):
		return f.generateCKOId("txn")
	case strings.Contains(fieldName, "customer_id") || strings.Contains(fieldName, "cus_id"):
		return f.generateCKOId("cus")
	case strings.Contains(fieldName, "source_id") || strings.Contains(fieldName, "src_id"):
		return f.generateCKOId("src")
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

// generateCKOId generates a Checkout.com style ID (prefix + base32 GUID)
// Format: prefix_base32guid (e.g., pay_uud7p4pbmtbulkno2fbax2apkm)
func (f *FakerAdapter) generateCKOId(prefix string) string {
	guid := uuid.New()
	base32Guid := f.base32EncodeGuid(guid[:])
	return prefix + "_" + base32Guid
}

// base32EncodeGuid encodes a UUID bytes to CKO-style base32 (26 chars)
func (f *FakerAdapter) base32EncodeGuid(data []byte) string {
	// CKO uses a custom base32 alphabet
	b32x := "abcdefghifklmnopqrstuvwxyz234567"

	// Reorder bytes like CKO does (similar to .NET GUID byte order)
	reordered := make([]byte, 16)
	// First 4 bytes reversed
	reordered[0] = data[3]
	reordered[1] = data[2]
	reordered[2] = data[1]
	reordered[3] = data[0]
	// Next 2 bytes reversed
	reordered[4] = data[5]
	reordered[5] = data[4]
	// Next 2 bytes reversed
	reordered[6] = data[7]
	reordered[7] = data[6]
	// Rest as-is
	copy(reordered[8:], data[8:])

	var dst strings.Builder

	// Process 5 bytes at a time (produces 8 base32 chars)
	for i := 0; i <= 10; i += 5 {
		dst.WriteByte(b32x[reordered[i]>>3])
		dst.WriteByte(b32x[((reordered[i]&0x07)<<2)|(reordered[i+1]>>6)])
		dst.WriteByte(b32x[(reordered[i+1]&0x3E)>>1])
		dst.WriteByte(b32x[((reordered[i+1]&0x01)<<4)|(reordered[i+2]>>4)])
		dst.WriteByte(b32x[((reordered[i+2]&0x0F)<<1)|(reordered[i+3]>>7)])
		dst.WriteByte(b32x[(reordered[i+3]&0x7C)>>2])
		dst.WriteByte(b32x[((reordered[i+3]&0x03)<<3)|((reordered[i+4]&0xE0)>>5)])
		dst.WriteByte(b32x[reordered[i+4]&0x1F])
	}

	// Handle last byte (produces 2 base32 chars)
	dst.WriteByte(b32x[reordered[15]>>3])
	dst.WriteByte(b32x[(reordered[15]&0x07)<<2])

	return dst.String()
}

// Helper to convert UUID to bytes in CKO order
func guidToBytes(guid uuid.UUID) []byte {
	bytes := make([]byte, 16)
	binary.LittleEndian.PutUint32(bytes[0:4], binary.BigEndian.Uint32(guid[0:4]))
	binary.LittleEndian.PutUint16(bytes[4:6], binary.BigEndian.Uint16(guid[4:6]))
	binary.LittleEndian.PutUint16(bytes[6:8], binary.BigEndian.Uint16(guid[6:8]))
	copy(bytes[8:], guid[8:])
	return bytes
}

// GenerateMultiple generates multiple values for a field
func (f *FakerAdapter) GenerateMultiple(fieldName, fieldType, format string, count int) []interface{} {
	results := make([]interface{}, count)
	for i := 0; i < count; i++ {
		results[i] = f.GenerateByType(fieldName, fieldType, format)
	}
	return results
}
