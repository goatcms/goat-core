package validator

import "github.com/goatcms/goatcore/messages"

const (
	// InvalidMinLength is key for too short strings
	InvalidMinLength = "length_min"
)

// MinStringValid add error message if string is shorten then some value
func MinStringValid(value string, basekey string, mm messages.MessageMap, min int) error {
	if len(value) < min {
		mm.Add(basekey, InvalidMinLength)
		return nil
	}
	return nil
}
