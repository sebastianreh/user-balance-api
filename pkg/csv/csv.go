package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
)

const (
	minRecords = 2
)

type CsvProcessor interface {
	ReadFile(file *multipart.FileHeader, recordValidator func(record []string) error) ([][]string, error)
}

type csvProcessor struct{}

func NewCsvProcessor() CsvProcessor {
	return &csvProcessor{}
}

func (c *csvProcessor) ReadFile(file *multipart.FileHeader, recordValidator func(record []string) error) ([][]string, error) {
	records := make([][]string, 0)

	src, err := file.Open()
	if err != nil {
		return records, err
	}
	defer src.Close()

	var fileBytes []byte
	fileBytes, err = io.ReadAll(src)
	if err != nil {
		return records, err
	}

	records, err = csvBytesToRecords(fileBytes, recordValidator)
	if err != nil {
		return records, err
	}

	return records, nil
}

func csvBytesToRecords(csvBytes []byte, validator func(record []string) error) ([][]string, error) {
	var records [][]string
	reader := csv.NewReader(bytes.NewReader(csvBytes))

	// Read the header and skip it to parse the records
	_, err := reader.Read()
	if err != nil {
		return records, err
	}

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() != "EOF" {
				return records, err
			}
			break
		}

		err = validator(record)
		if err != nil {
			return records, err
		}

		records = append(records, record)
	}

	if len(records) < minRecords {
		return records, errors.New("error reading csv bytes - records < 2")
	}

	return records, nil
}
