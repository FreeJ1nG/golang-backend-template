package cms

import (
	"fmt"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5"
)

type utils struct {
}

func NewUtil() *utils {
	return &utils{}
}

func (u *utils) ConvertColumnsToAttributes(columns []models.Column) (attributes []string) {
	attributes = make([]string, 0)
	for _, column := range columns {
		if column.ValueType != "data_type" {
			continue
		}
		attributes = append(attributes, column.ColumnName)
	}
	return
}

func (u *utils) ValidateData(data map[string]interface{}, columns []models.Column) (err error) {
	for key, value := range data {
		for _, column := range columns {
			if column.ColumnName == key {
				err = column.ValidateData(value)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

func (u *utils) ConvertRowToMap(row pgx.Row, attributes []string) (res map[string]interface{}, err error) {
	scanTargets := make([]interface{}, len(attributes))
	for i := range attributes {
		var result interface{}
		scanTargets[i] = &result
	}

	err = row.Scan(scanTargets...)
	if err != nil {
		return
	}

	fmt.Println(scanTargets)

	res = make(map[string]interface{})
	for i, attribute := range attributes {
		res[attribute] = scanTargets[i]
	}
	return
}

func (u *utils) ConvertAttributesToPsqlAttributes(attributes []string) (res []string) {
	res = make([]string, len(attributes))
	for i, attribute := range attributes {
		res[i] = strcase.ToSnake(attribute)
	}
	return
}
