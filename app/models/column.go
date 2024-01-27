package models

type Column struct {
	ValueType  string `json:"valueType"`
	ColumnName string `json:"columnName"`
	Value      string `json:"value"`
}
