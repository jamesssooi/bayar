package bayar

import (
	"fmt"
	"time"

	"google.golang.org/api/sheets/v4"
)

// Expense is a struct describing a single expense item.
type Expense struct {
	Label    string  `json:"label"`
	Category string  `json:"category"`
	Cost     float32 `json:"cost"`
}

// NewExpense creates a new Expense instance.
func NewExpense(label string, category string, cost float32) Expense {
	e := Expense{
		Label:    label,
		Category: category,
		Cost:     cost,
	}
	return e
}

func (e Expense) insertIntoSpreadsheet(spreadsheetID string, sheet string) (int, error) {
	client, err := getGoogleClient()
	if err != nil {
		return -1, err
	}

	service, err := sheets.New(client)
	if err != nil {
		return -1, err
	}

	var vr sheets.ValueRange
	vals := []interface{}{
		nil,
		time.Now().Local().Format("2019-01-01"),
		e.Label,
		e.Category,
		nil,
		nil,
		nil,
		e.Cost,
	}
	vr.Values = append(vr.Values, vals)

	row, err := getFirstEmptyRow(service, spreadsheetID, sheet, "B", 4)
	if err != nil {
		return -1, err
	}

	insertRange := fmt.Sprintf("%s!A%d:H%d", sheet, row, row)
	_, insertErr := service.Spreadsheets.Values.Update(spreadsheetID, insertRange, &vr).ValueInputOption("USER_ENTERED").Do()
	if insertErr != nil {
		return -1, insertErr
	}

	return row, nil
}

// Scans a specific range of values and returns the first empty row.
func getFirstEmptyRow(service *sheets.Service, spreadsheetID string, sheet string, column string, offset int) (int, error) {
	scanRange := fmt.Sprintf("%s!%s%d:%s", sheet, column, offset, column)
	response, err := service.Spreadsheets.Values.Get(spreadsheetID, scanRange).Do()
	if err != nil {
		return -1, err
	}
	emptyrow := len(response.Values) + offset
	return emptyrow, nil
}
