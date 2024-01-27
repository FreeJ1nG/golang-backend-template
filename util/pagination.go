package util

import (
	"math"

	"github.com/FreeJ1nG/backend-template/app/models"
)

type paginator struct {
}

func NewPaginator() *paginator {
	return &paginator{}
}

func (p *paginator) Paginate(length int, opts models.PaginationOptions) (res models.PaginationMetadata) {
	res.TotalPages = int(math.Ceil(float64(length) / float64(opts.ItemPerPage)))
	res.LowerBound = (opts.CurrentPage-1)*opts.ItemPerPage + 1
	res.UpperBound = opts.CurrentPage * opts.ItemPerPage
	return
}
