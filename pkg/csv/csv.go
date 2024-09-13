package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
)

const (
	minRecords = 2
)

type CsvProcessor struct{}

func NewCsvProcessor() CsvProcessor {
	return CsvProcessor{}
}

func (*CsvProcessor) CsvBytesToRecords(csvBytes []byte, validator func(record []string) error) ([][]string, error) {
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
