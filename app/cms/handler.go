package cms

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/dto"
	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/pagination"
	"github.com/FreeJ1nG/backend-template/util"
	"github.com/gorilla/mux"
	"github.com/iancoleman/strcase"
)

type handler struct {
	cmsService interfaces.CmsService
}

func NewHandler(cmsService interfaces.CmsService) *handler {
	return &handler{
		cmsService: cmsService,
	}
}

func (h *handler) GetTableInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["tableName"]

	res, status, err := h.cmsService.GetTableInfo(tableName)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}

	util.EncodeSuccessResponse(w, res, status, nil)
}

func (h *handler) GetTableData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["tableName"]

	relatedTables := r.URL.Query()["relatedTables[]"]

	opts := util.ParseRequestBody[pagination.Options](w, r)
	res, metadata, status, err := h.cmsService.GetTableData(tableName, relatedTables, &opts)

	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}

	util.EncodeSuccessResponse(w, res, status, &metadata)
}

func (h *handler) CreateTableData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["tableName"]

	body := util.ParseRequestBody[map[string]interface{}](w, r)
	body = util.ConvertMapKeys(body, strcase.ToSnake)
	res, status, err := h.cmsService.CreateTableData(tableName, body)
	res = util.ConvertMapKeys(res, strcase.ToLowerCamel)

	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}

	util.EncodeSuccessResponse(w, res, status, nil)
}

func (h *handler) UpdateTableData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["tableName"]
	pkString := vars["pk"]

	body := util.ParseRequestBody[map[string]interface{}](w, r)
	body = util.ConvertMapKeys(body, strcase.ToSnake)

	pk, status, err := util.ParseStringToInt(pkString)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}

	res, status, err := h.cmsService.UpdateTableData(tableName, pk, body)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}

	util.EncodeSuccessResponse(w, res, status, nil)
}

func (h *handler) DeleteTableData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableName := vars["tableName"]
	pkString := vars["pk"]

	pk, status, err := util.ParseStringToInt(pkString)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}

	status, err = h.cmsService.DeleteTableData(tableName, pk)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}

	util.EncodeSuccessResponse(w, dto.EmptySuccessMessage, status, nil)
}
