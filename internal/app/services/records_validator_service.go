package services

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	customStr "github.com/sebastianreh/user-balance-api/pkg/strings"
)

const (
	minRecordLen = 4
)

func recordValidator(record []string) error {
	validators := []func(string) error{
		validateID, validateUserID, validateAmount, validateDatetime,
	}

	if err := validateRecordLength(record); err != nil {
		return err
	}

	for i, validator := range validators {
		if err := validator(record[i]); err != nil {
			return err
		}
	}

	return nil
}

func validateRecordLength(record []string) error {
	if len(record) < minRecordLen {
		return errors.New("record fields are below required")
	}
	return nil
}

func validateID(id string) error {
	return validateIntValue("id", id)
}

func validateUserID(userID string) error {
	return validateIntValue("userID", userID)
}

func validateAmount(amount string) error {
	if customStr.IsEmpty(amount) {
		return errors.New("amount field is empty")
	}

	if _, err := strconv.ParseFloat(amount, 64); err != nil {
		return errors.New("amount field is not a valid float")
	}

	return nil
}

func validateIntValue(fieldName, value string) error {
	if customStr.IsEmpty(value) {
		return fmt.Errorf("%s field is empty", fieldName)
	}

	if _, err := strconv.Atoi(value); err != nil {
		return fmt.Errorf("%s field is empty", fieldName)
	}

	return nil
}

func validateDatetime(datetime string) error {
	if customStr.IsEmpty(datetime) {
		return errors.New("datetime field is empty")
	}

	if _, err := time.Parse(time.RFC3339, datetime); err != nil {
		return errors.New("datetime field is not in valid ISO 8601 format")
	}

	return nil
}
