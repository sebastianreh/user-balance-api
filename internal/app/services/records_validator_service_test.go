package services

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RecordValidator(t *testing.T) {
	t.Run("When record is valid", func(t *testing.T) {
		record := []string{"1", "123", "100.50", "2024-09-13T10:00:00Z"}
		err := recordValidator(record)
		assert.Nil(t, err)
	})

	t.Run("When record has invalid length", func(t *testing.T) {
		record := []string{"1", "123", "100.50"}
		err := recordValidator(record)
		assert.NotNil(t, err)
		assert.Equal(t, "record fields are below required", err.Error())
	})

	t.Run("When ID field is invalid", func(t *testing.T) {
		record := []string{"abc", "123", "100.50", "2024-09-13T10:00:00Z"}
		err := recordValidator(record)
		assert.NotNil(t, err)
		assert.Equal(t, "id field is empty", err.Error())
	})

	t.Run("When UserID field is invalid", func(t *testing.T) {
		record := []string{"1", "abc", "100.50", "2024-09-13T10:00:00Z"}
		err := recordValidator(record)
		assert.NotNil(t, err)
		assert.Equal(t, "userID field is empty", err.Error())
	})

	t.Run("When Amount field is empty", func(t *testing.T) {
		record := []string{"1", "123", "", "2024-09-13T10:00:00Z"}
		err := recordValidator(record)
		assert.NotNil(t, err)
		assert.Equal(t, "amount field is empty", err.Error())
	})

	t.Run("When Amount field is not a valid float", func(t *testing.T) {
		record := []string{"1", "123", "abc", "2024-09-13T10:00:00Z"}
		err := recordValidator(record)
		assert.NotNil(t, err)
		assert.Equal(t, "amount field is not a valid float", err.Error())
	})

	t.Run("When Datetime field is empty", func(t *testing.T) {
		record := []string{"1", "123", "100.50", ""}
		err := recordValidator(record)
		assert.NotNil(t, err)
		assert.Equal(t, "datetime field is empty", err.Error())
	})

	t.Run("When Datetime field is not valid ISO 8601", func(t *testing.T) {
		record := []string{"1", "123", "100.50", "13/09/2024"}
		err := recordValidator(record)
		assert.NotNil(t, err)
		assert.Equal(t, "datetime field is not in valid ISO 8601 format", err.Error())
	})
}

func Test_ValidateID(t *testing.T) {
	t.Run("When ID is valid", func(t *testing.T) {
		err := validateID("123")
		assert.Nil(t, err)
	})

	t.Run("When ID is empty", func(t *testing.T) {
		err := validateID("")
		assert.NotNil(t, err)
		assert.Equal(t, "id field is empty", err.Error())
	})

	t.Run("When ID is not an integer", func(t *testing.T) {
		err := validateID("abc")
		assert.NotNil(t, err)
		assert.Equal(t, "id field is empty", err.Error())
	})
}

func Test_ValidateUserID(t *testing.T) {
	t.Run("When UserID is valid", func(t *testing.T) {
		err := validateUserID("456")
		assert.Nil(t, err)
	})

	t.Run("When UserID is empty", func(t *testing.T) {
		err := validateUserID("")
		assert.NotNil(t, err)
		assert.Equal(t, "userID field is empty", err.Error())
	})

	t.Run("When UserID is not an integer", func(t *testing.T) {
		err := validateUserID("xyz")
		assert.NotNil(t, err)
		assert.Equal(t, "userID field is empty", err.Error())
	})
}

func Test_ValidateAmount(t *testing.T) {
	t.Run("When Amount is valid", func(t *testing.T) {
		err := validateAmount("100.50")
		assert.Nil(t, err)
	})

	t.Run("When Amount is empty", func(t *testing.T) {
		err := validateAmount("")
		assert.NotNil(t, err)
		assert.Equal(t, "amount field is empty", err.Error())
	})

	t.Run("When Amount is not a valid float", func(t *testing.T) {
		err := validateAmount("abc")
		assert.NotNil(t, err)
		assert.Equal(t, "amount field is not a valid float", err.Error())
	})
}

func Test_ValidateDatetime(t *testing.T) {
	t.Run("When Datetime is valid", func(t *testing.T) {
		err := validateDatetime("2024-09-13T10:00:00Z")
		assert.Nil(t, err)
	})

	t.Run("When Datetime is empty", func(t *testing.T) {
		err := validateDatetime("")
		assert.NotNil(t, err)
		assert.Equal(t, "datetime field is empty", err.Error())
	})

	t.Run("When Datetime is not valid ISO 8601", func(t *testing.T) {
		err := validateDatetime("13/09/2024")
		assert.NotNil(t, err)
		assert.Equal(t, "datetime field is not in valid ISO 8601 format", err.Error())
	})
}
