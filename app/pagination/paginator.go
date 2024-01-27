package pagination

import (
	"context"
	"fmt"
	"math"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Options struct {
	ItemPerPage int `json:"itemPerPage"`
	CurrentPage int `json:"currentPage"`
}

type Metadata struct {
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type Paginator struct {
	mainDB *pgxpool.Pool
}

func NewPaginator(mainDB *pgxpool.Pool) *Paginator {
	return &Paginator{
		mainDB: mainDB,
	}
}

func (p *Paginator) getTotalDataInstance(tableName string) (res int, err error) {
	ctx := context.Background()
	var result []int
	err = pgxscan.Select(
		ctx,
		p.mainDB,
		&result,
		fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName),
	)
	if len(result) != 1 {
		err = fmt.Errorf("unable to get total data count")
		return
	}
	res = result[0]
	return
}

func (p *Paginator) GetPaginationInfo(tableName string, opts *Options) (offset int, limit int, metadata Metadata, err error) {
	offset = opts.ItemPerPage * (opts.CurrentPage - 1)
	limit = opts.ItemPerPage
	metadata.TotalItems, err = p.getTotalDataInstance(tableName)
	if err != nil {
		return
	}
	metadata.TotalPages = int(math.Ceil(float64(metadata.TotalItems) / float64(opts.ItemPerPage)))
	return
}
