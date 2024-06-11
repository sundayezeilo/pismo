package validators

import (
	"errors"
	"regexp"
)

func ValidateCreateAccountReq(documentNumber string) error {
	pattern := `^\d{11}$`
	re := regexp.MustCompile(pattern)

	if documentNumber == "" {
		return errors.New("document_number is required")
	}

	if !re.MatchString(documentNumber) {
		return errors.New("invalid document_number")
	}

	return nil
}
