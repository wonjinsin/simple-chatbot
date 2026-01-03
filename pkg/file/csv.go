package file

import (
	"encoding/csv"
	"fmt"
	"os"
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
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	headers := records[0]
	if len(headers) == 0 {
		return nil, fmt.Errorf("CSV file has no columns")
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
