package csvutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// CSV is a utility class for working with CSV files.
type CSV struct {
	file     *os.File
	reader   *csv.Reader
	rowNames map[string][]string
}

// NewCSV creates a new CSV object from a file path.
func NewCSV(path string) (*CSV, error) {
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// Create a new CSV reader
	reader := csv.NewReader(file)

	value, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	rowNames := make(map[string][]string)

	for _, v := range value {
		rowNames[v[0]] = v
	}

	// Create and return a new CSV object
	return &CSV{
		file:     file,
		reader:   reader,
		rowNames: rowNames,
	}, nil
}

// Close closes the CSV file.
func (c *CSV) Close() error {
	return c.file.Close()
}

// GetRows returns all rows from the CSV file.
func (c *CSV) GetRows() ([][]string, error) {
	return c.reader.ReadAll()
}

func (c *CSV) GetRowNames() map[string][]string {
	return c.rowNames
}

func (c *CSV) GetRowIdx(name string) []string {
	return c.rowNames[name]
}

// GetRow returns a specific row from the CSV file.
func (c *CSV) GetRow(index int) ([]string, error) {
	// Read rows until we get to the desired index
	for i := 0; i < index; i++ {
		if _, err := c.reader.Read(); err != nil {
			return nil, err
		}
	}

	// Read the desired row and return it
	return c.reader.Read()
}

// GetCellValue returns the value of a specific cell in the CSV file.
func (c *CSV) GetCellValue(row int, col int) (string, error) {
	// Get the desired row
	rowData, err := c.GetRow(row)
	if err != nil {
		return "", err
	}

	// Get the desired cell value
	if col >= len(rowData) {
		return "", nil
	}
	return rowData[col], nil
}

// GetRowsByColumn returns all rows from the CSV file that match the specified column name.
func (c *CSV) GetRowsByColumn(columnName string) ([][]string, error) {
	// Read the header row to get the column indices
	header, err := c.reader.Read()
	if err != nil {
		return nil, err
	}

	// Find the index of the desired column
	colIndex := -1
	for i, name := range header {
		if name == columnName {
			colIndex = i
			break
		}
	}

	if colIndex == -1 {
		return nil, fmt.Errorf("column not found: %s", columnName)
	}

	// Read rows until we find all the rows that match the column name
	var matchingRows [][]string
	for {
		row, err := c.reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(row) > colIndex && row[colIndex] != "" {
			matchingRows = append(matchingRows, row)
		}
	}

	return matchingRows, nil
}

// GetValueByField returns the value of a specific field in the CSV file.
func (c *CSV) GetValueByField(field string) (string, error) {
	// Read the header row to get the field indices
	header, err := c.reader.Read()
	if err != nil {
		return "", err
	}

	// Find the index of the desired field
	fieldIndex := -1
	for i, name := range header {
		if name == field {
			fieldIndex = i
			break
		}
	}

	if fieldIndex == -1 {
		return "", fmt.Errorf("field not found: %s", field)
	}

	// Read rows until we find the row with the desired field
	for {
		row, err := c.reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		if row[0] == field {
			return row[fieldIndex], nil
		}
	}

	return "", fmt.Errorf("field not found: %s", field)
}
