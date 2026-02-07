package action

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

// ExportCSV creates a CSV export action.
// It exports the selected models to a CSV file.
func ExportCSV(filename string) *BaseAction {
	return New("Export as CSV").
		SetIcon("download").
		Handle(func(ctx *ActionContext) error {
			return exportToCSV(ctx.Models, filename)
		})
}

// Delete creates a bulk delete action.
// It deletes all selected models from the database.
func Delete() *BaseAction {
	return New("Delete").
		SetIcon("trash").
		Destructive().
		Confirm("Are you sure you want to delete these items?").
		ConfirmButton("Delete").
		Handle(func(ctx *ActionContext) error {
			for _, model := range ctx.Models {
				if err := ctx.DB.Delete(model).Error; err != nil {
					return err
				}
			}
			return nil
		})
}

// Approve creates an approval action.
// It can be customized with a field name to update.
func Approve() *BaseAction {
	return New("Approve").
		SetIcon("check").
		Confirm("Are you sure you want to approve these items?").
		Handle(func(ctx *ActionContext) error {
			// Default approval logic - can be overridden
			for _, model := range ctx.Models {
				// Try to set a common "status" or "approved" field
				v := reflect.ValueOf(model)
				if v.Kind() == reflect.Ptr {
					v = v.Elem()
				}

				if v.Kind() == reflect.Struct {
					// Try to find and set status field
					statusField := v.FieldByName("Status")
					if statusField.IsValid() && statusField.CanSet() {
						statusField.SetString("approved")
						ctx.DB.Save(model)
						continue
					}

					// Try to find and set approved field
					approvedField := v.FieldByName("Approved")
					if approvedField.IsValid() && approvedField.CanSet() {
						approvedField.SetBool(true)
						ctx.DB.Save(model)
						continue
					}
				}
			}
			return nil
		})
}

// exportToCSV exports models to a CSV file
func exportToCSV(models []interface{}, filename string) error {
	if len(models) == 0 {
		return fmt.Errorf("no models to export")
	}

	// Create exports directory if it doesn't exist
	exportsDir := filepath.Join("storage", "exports")
	if err := os.MkdirAll(exportsDir, 0755); err != nil {
		return fmt.Errorf("failed to create exports directory: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	if filename == "" {
		filename = fmt.Sprintf("export_%s.csv", timestamp)
	} else {
		ext := filepath.Ext(filename)
		name := filename[:len(filename)-len(ext)]
		filename = fmt.Sprintf("%s_%s%s", name, timestamp, ext)
	}

	filePath := filepath.Join(exportsDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Get field names from the first model
	firstModel := models[0]
	v := reflect.ValueOf(firstModel)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("model must be a struct")
	}

	t := v.Type()
	var headers []string
	var fieldIndices []int

	// Collect field names and indices
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		headers = append(headers, field.Name)
		fieldIndices = append(fieldIndices, i)
	}

	// Write headers
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV headers: %w", err)
	}

	// Write data rows
	for _, model := range models {
		v := reflect.ValueOf(model)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		var row []string
		for _, idx := range fieldIndices {
			fieldValue := v.Field(idx)
			row = append(row, fmt.Sprintf("%v", fieldValue.Interface()))
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	return nil
}
