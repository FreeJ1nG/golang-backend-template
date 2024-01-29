package interfaces

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/FreeJ1nG/backend-template/app/pagination"
	"github.com/jackc/pgx/v5"
)

type CmsRepository interface {
	GetTableDataTypes(tableName string, columnConverter *func(s string) string) (columns []models.Column, err error)
	GetTableData(tableName string, opts *pagination.Options) (res []map[string]interface{}, metadata pagination.Metadata, err error)
	CreateTableData(tableName string, data map[string]interface{}, attributes []string) (res map[string]interface{}, err error)
	UpdateTableDataByPk(tableName string, pk int, data map[string]interface{}, attributes []string) (res map[string]interface{}, err error)
	DeleteTableDataByPk(tableName string, pk int) (err error)
}

type CmsService interface {
	GetTableInfo(tableName string) (res []models.Column, status int, err error)
	GetTableData(tableName string, opts *pagination.Options) (res []map[string]interface{}, metadata pagination.Metadata, status int, err error)
	CreateTableData(tableName string, data map[string]interface{}) (res map[string]interface{}, status int, err error)
	UpdateTableData(tableName string, pk int, data map[string]interface{}) (res map[string]interface{}, status int, err error)
	DeleteTableData(tableName string, pk int) (status int, err error)
}

type CmsHandler interface {
	GetTableInfo(w http.ResponseWriter, r *http.Request)
	GetTableData(w http.ResponseWriter, r *http.Request)
	CreateTableData(w http.ResponseWriter, r *http.Request)
	UpdateTableData(w http.ResponseWriter, r *http.Request)
	DeleteTableData(w http.ResponseWriter, r *http.Request)
}

type CmsUtil interface {
	ConvertColumnsToAttributes(columns []models.Column) (attributes []string)
	ValidateData(data map[string]interface{}, columns []models.Column) (err error)
	ConvertRowToMap(row pgx.Row, attributes []string) (res map[string]interface{}, err error)
	ConvertAttributesToPsqlAttributes(attributes []string) (res []string)
}
