package cms

import (
	"context"
	"fmt"
	"strings"

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

func (r *repository) CreateTableData(tableName string, data map[string]interface{}) (res map[string]interface{}, err error) {
	ctx := context.Background()

	var columns []string
	var placeholders []string
	var values []interface{}
	i := 1

	for key, value := range data {
		columns = append(columns, key)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		values = append(values, value)
		i++
	}

	joinedColumns := strings.Join(columns, ",")
	joinedPlaceholders := strings.Join(placeholders, ",")
	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING %s;",
		tableName,
		joinedColumns,
		joinedPlaceholders,
		"id,"+joinedColumns,
	)

	row := r.mainDB.QueryRow(
		ctx,
		sql,
		values...,
	)

	scanTargets := make([]interface{}, len(columns)+1)
	for i := range scanTargets {
		var result interface{}
		scanTargets[i] = &result
	}

	err = row.Scan(scanTargets...)
	if err != nil {
		return
	}

	res = make(map[string]interface{})
	for i, value := range scanTargets[1:] {
		res[columns[i]] = value
	}
	return
}
