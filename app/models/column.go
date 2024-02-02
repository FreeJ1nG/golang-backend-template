package models

import (
	"fmt"
	"strconv"
	"strings"
)

type Column struct {
	ValueType  string `json:"valueType"`
	ColumnName string `json:"columnName"`
	Value      string `json:"value"`
}

func (c *Column) validateInt(data int) (err error) {
	if c.Value != "integer" {
		err = fmt.Errorf("column %s has invalid value, expected %s, got %s", c.ColumnName, c.Value, "integer")
		return
	}
	return
}

func (c *Column) validateString(data string) (err error) {
	vSplit := strings.Split(c.Value, " ")
	isVarChar := (len(vSplit) == 2 && vSplit[0] == "character" && strings.Contains(vSplit[1], "varying") && len(vSplit[1]) >= 9)
	if c.Value != "text" && !isVarChar {
		return fmt.Errorf("column %s has invalid value, expected %s, got %s", c.ColumnName, c.Value, "string")
	}
	if isVarChar {
		varyingStr := vSplit[1]
		maxLength, err := strconv.Atoi(varyingStr[8 : len(varyingStr)-1])
		if err != nil {
			return fmt.Errorf("unable to convert max length to integer")
		}
		if len(data) > maxLength {
			return fmt.Errorf("expected a string with max length of %d, got %d", maxLength, len(data))
		}
	}
	return
}

func (c *Column) validateBoolean(data bool) (err error) {
	if c.Value != "boolean" {
		err = fmt.Errorf("column %s has invalid value, expected %s, got %s", c.ColumnName, c.Value, "boolean")
		return
	}
	return
}

func (c *Column) ValidateData(data interface{}) (err error) {
	if c.ValueType == "data_type" {
		switch v := data.(type) {
		case int:
			err = c.validateInt(v)
		case string:
			err = c.validateString(v)
		case bool:
			err = c.validateBoolean(v)
		default:
			err = fmt.Errorf("this data type is currently unsupported")
		}
	}
	return
}
