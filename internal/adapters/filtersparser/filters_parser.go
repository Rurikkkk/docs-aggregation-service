package filtersparser

import (
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/extrame/xls"
)

type FiltersParser struct{}

func NewFiltersParser() *FiltersParser {
	return &FiltersParser{}
}

func (fr *FiltersParser) ParseFilters(reader io.ReadSeeker) ([]string, error) {
	file, err := xls.OpenReader(reader, "utf-8")
	if err != nil {
		log.Printf("[FiltersParser] Opening file from reader failed: %v", err)
		return nil, err
	}

	if file.NumSheets() == 0 {
		log.Printf("[FiltersParser] File is empty")
		return []string{}, nil
	}

	sheet := file.GetSheet(0)
	var fiscalDriveNumbers []string
	for i := 0; i <= int(sheet.MaxRow); i++ {
		row := sheet.Row(i)
		if row == nil {
			continue
		}
		value := strings.TrimSpace(row.Col(0))
		if value == "" {
			continue
		}
		_, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("[FiltersParser] Value is not a number, skipping: %s", value)
			continue
		}
		fiscalDriveNumbers = append(fiscalDriveNumbers, value)
	}
	return fiscalDriveNumbers, nil
}
