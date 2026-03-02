package filtersrepo

import (
	"log"
	"strings"

	"github.com/extrame/xls"
)

type FiltersRepo struct {
	path string
}

func NewFiltersRepo(path string) *FiltersRepo {
	return &FiltersRepo{path: path}
}

func (fr *FiltersRepo) GetFilters() ([]string, error) {
	file, err := xls.Open(fr.path, "utf-8")
	if err != nil {
		log.Printf("[FiltersRepo] Opening %s failed: %v", fr.path, err)
		return nil, err
	}

	if file.NumSheets() == 0 {
		log.Printf("[FiltersRepo] File is empty")
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
		fiscalDriveNumbers = append(fiscalDriveNumbers, value)
	}
	return fiscalDriveNumbers, nil
}
