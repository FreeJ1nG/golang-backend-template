package cms

import (
	"context"
	"fmt"

	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/FreeJ1nG/backend-template/app/pagination"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	mainDB    *pgxpool.Pool
	paginator *pagination.Paginator
}

func NewRepository(mainDB *pgxpool.Pool, paginator *pagination.Paginator) *repository {
	return &repository{
		mainDB:    mainDB,
		paginator: paginator,
	}
}

func (r *repository) GetTableDataTypes(tableName string) (columns []models.Column, err error) {
	ctx := context.Background()
	err = pgxscan.Select(
		ctx,
		r.mainDB,
		&columns,
		`
		SELECT 'data_type' AS value_type,
		attname AS column_name,
		format_type(atttypid, atttypmod) AS value
		FROM pg_attribute
		WHERE attrelid = $1::regclass
		AND attnum > 0
		AND NOT attisdropped
		UNION ALL
		SELECT 'constraint' AS info,
		a.attname AS column_name,
		'is_primary_key' AS value
		FROM pg_constraint c
		JOIN pg_attribute a ON a.attnum = ANY(c.conkey) AND a.attrelid = c.conrelid
		WHERE contype = 'p'
		AND conrelid = $1::regclass
		UNION ALL
		SELECT 'constraint' AS info,
		a.attname AS column_name,
		'is_foreign_key' AS value
		FROM pg_constraint c
		JOIN pg_attribute a ON a.attnum = ANY(c.conkey) AND a.attrelid = c.conrelid
		WHERE contype = 'f'
		AND conrelid = $1::regclass
		UNION ALL
		SELECT 'constraint' AS info,
		conname AS column_name,
		'is_unique' AS value
		FROM pg_constraint
		WHERE contype = 'u'
		AND conrelid = $1::regclass;`,
		tableName,
	)
	if err != nil {
		return
	}
	return
}

func (r *repository) GetTableData(tableName string, opts *pagination.Options) (res []map[string]interface{}, metadata pagination.Metadata, err error) {
	ctx := context.Background()
	offset, limit, metadata, err := r.paginator.Paginate(tableName, opts)
	if err != nil {
		return
	}
	err = pgxscan.Select(
		ctx,
		r.mainDB,
		&res,
		fmt.Sprintf("SELECT * FROM %s OFFSET %d LIMIT %d", tableName, offset, limit),
	)
	return
}
