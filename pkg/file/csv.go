package file

import (
	"encoding/csv"
	"os"

	"github.com/wonjinsin/simple-chatbot/internal/constants"
	"github.com/wonjinsin/simple-chatbot/pkg/errors"
)

// ReadCSVToMapArray reads a CSV file and converts it to an array of maps.
// The first row is treated as column headers, and subsequent rows are data.
// Each row is converted to a map where keys are column names and values are cell values.
//
// Parameters:
//   - filePath: path to the CSV file
//
// Returns:
//   - []map[string]string: array of maps, each representing a row
//   - error: any error encountered during file reading or parsing
//
// Behavior:
//   - If a row has fewer columns than the header, missing values are filled with empty strings
//   - If a row has more columns than the header, extra columns are ignored
//   - Empty rows are automatically skipped by the csv.Reader
func ReadCSVToMapArray(filePath string) ([]map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open CSV file")
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read CSV file")
	}

	if len(records) == 0 {
		return nil, errors.New(constants.InvalidParameter, "CSV file is empty", nil)
	}

	headers := records[0]
	if len(headers) == 0 {
		return nil, errors.New(constants.InvalidParameter, "CSV file has no columns", nil)
	}

	result := make([]map[string]string, 0, len(records)-1)

	for i := 1; i < len(records); i++ {
		row := records[i]
		rowMap := make(map[string]string, len(headers))

		for j, header := range headers {
			if j < len(row) {
				rowMap[header] = row[j]
			} else {
				rowMap[header] = ""
			}
		}

		result = append(result, rowMap)
	}

	return result, nil
}
