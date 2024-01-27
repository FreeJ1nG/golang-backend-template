package cms

import (
	"fmt"
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/FreeJ1nG/backend-template/app/pagination"
)

type service struct {
	cmsRepo interfaces.CmsRepository
}

func NewService(cmsRepo interfaces.CmsRepository) *service {
	return &service{
		cmsRepo: cmsRepo,
	}
}

func (s *service) GetTableInfo(tableName string) (res []models.Column, status int, err error) {
	status = http.StatusOK
	res, err = s.cmsRepo.GetTableDataTypes(tableName)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to get table data type: %s", err.Error())
		return
	}
	return
}

func (s *service) GetTableData(tableName string, opts *pagination.Options) (res []map[string]interface{}, metadata pagination.Metadata, status int, err error) {
	status = http.StatusOK
	res, metadata, err = s.cmsRepo.GetTableData(tableName, opts)
	if err != nil {
		status = http.StatusInternalServerError
		err = fmt.Errorf("unable to get table data: %s", err.Error())
		return
	}
	return
}
