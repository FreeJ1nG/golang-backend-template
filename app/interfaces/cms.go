package interfaces

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/FreeJ1nG/backend-template/app/pagination"
)

type CmsRepository interface {
	GetTableDataTypes(tableName string) (columns []models.Column, err error)
	GetTableData(tableName string, opts *pagination.Options) (res []map[string]interface{}, metadata pagination.Metadata, err error)
	CreateTableData(tableName string, data map[string]interface{}) (res map[string]interface{}, err error)
}

type CmsService interface {
	GetTableInfo(tableName string) (res []models.Column, status int, err error)
	GetTableData(tableName string, opts *pagination.Options) (res []map[string]interface{}, metadata pagination.Metadata, status int, err error)
	CreateTableData(tableName string, data map[string]interface{}) (res map[string]interface{}, status int, err error)
}

type CmsHandler interface {
	GetTableInfo(w http.ResponseWriter, r *http.Request)
	GetTableData(w http.ResponseWriter, r *http.Request)
	CreateTableData(w http.ResponseWriter, r *http.Request)
}
