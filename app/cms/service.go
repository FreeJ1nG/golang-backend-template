package cms

import (
	"fmt"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/FreeJ1nG/backend-template/app/pagination"
	"github.com/iancoleman/strcase"
)

type service struct {
	cmsRepo interfaces.CmsRepository
	cmsUtil interfaces.CmsUtil
}

func NewService(cmsRepo interfaces.CmsRepository, cmsUtil interfaces.CmsUtil) *service {
	return &service{
		cmsRepo: cmsRepo,
		cmsUtil: cmsUtil,
	}
}

func (s *service) GetTableInfo(tableName string) (res []models.Column, status int, err error) {
	status = http.StatusOK

	converter := strcase.ToLowerCamel
	res, err = s.cmsRepo.GetTableDataTypes(tableName, &converter)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to get table data type: %s", err.Error())
		return
	}

	return
}

func (s *service) GetTableData(tableName string, relatedTables []string, opts *pagination.Options) (res map[string]interface{}, metadata pagination.Metadata, status int, err error) {
	status = http.StatusOK

	res = make(map[string]interface{})
	var tableData []map[string]interface{}
	tableData, metadata, err = s.cmsRepo.GetTableData(tableName, opts)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to get table data: %s", err.Error())
		return
	}
	res["tableData"] = tableData

	for _, relatedTableName := range relatedTables {
		var relatedTableData []map[string]interface{}
		relatedTableData, err = s.cmsRepo.GetRelatedTableData(tableName, relatedTableName)
		if err != nil {
			err = nil
			continue
		}
		res[relatedTableName] = relatedTableData
	}

	return
}

func (s *service) CreateTableData(tableName string, data map[string]interface{}) (res map[string]interface{}, status int, err error) {
	status = http.StatusOK

	converter := strcase.ToLowerCamel
	columns, err := s.cmsRepo.GetTableDataTypes(tableName, &converter)
	if err != nil {
		status = http.StatusBadRequest
		err = fmt.Errorf("unable to get table data types: %s", err.Error())
		return
	}

	err = s.cmsUtil.ValidateData(data, columns)
	if err != nil {
		status = http.StatusBadRequest
		err = fmt.Errorf("invalid data: %s", err.Error())
		return
	}

	tableAttributes := s.cmsUtil.ConvertColumnsToAttributes(columns)
	res, err = s.cmsRepo.CreateTableData(tableName, data, tableAttributes)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to create table data: %s", err.Error())
		return
	}

	return
}

func (s *service) UpdateTableData(tableName string, pk int, data map[string]interface{}) (res map[string]interface{}, status int, err error) {
	status = http.StatusOK

	converter := strcase.ToLowerCamel
	columns, err := s.cmsRepo.GetTableDataTypes(tableName, &converter)
	if err != nil {
		status = http.StatusBadRequest
		err = fmt.Errorf("unable to get table data type: %s", err.Error())
		return
	}

	err = s.cmsUtil.ValidateData(data, columns)
	if err != nil {
		status = http.StatusBadRequest
		err = fmt.Errorf("unable to validate data type: %s", err.Error())
		return
	}

	tableAttributes := s.cmsUtil.ConvertColumnsToAttributes(columns)
	res, err = s.cmsRepo.UpdateTableDataByPk(tableName, pk, data, tableAttributes)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to update table data with primary key %d: %s", pk, err.Error())
		return
	}

	return
}

func (s *service) DeleteTableData(tableName string, pk int) (status int, err error) {
	status = http.StatusOK

	err = s.cmsRepo.DeleteTableDataByPk(tableName, pk)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to delete table data with primary key %d: %s", pk, err.Error())
		return
	}

	return
}
