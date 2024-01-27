package cms

import (
	"net/http"

	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/pagination"
	"github.com/FreeJ1nG/backend-template/util"
	"github.com/gorilla/mux"
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
	opts := util.ParseRequestBody[pagination.Options](w, r)
	res, metadata, status, err := h.cmsService.GetTableData(tableName, &opts)
	if err != nil {
		util.EncodeErrorResponse(w, err.Error(), status)
		return
	}
	util.EncodeSuccessResponse(w, res, status, &metadata)
}
