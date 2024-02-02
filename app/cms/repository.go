package cms

import (
	"context"
	"fmt"
	"strings"

	"github.com/FreeJ1nG/backend-template/app/interfaces"
	"github.com/FreeJ1nG/backend-template/app/models"
	"github.com/FreeJ1nG/backend-template/app/pagination"
	"github.com/FreeJ1nG/backend-template/util"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/iancoleman/strcase"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	mainDB    *pgxpool.Pool
	cmsUtil   interfaces.CmsUtil
	paginator *pagination.Paginator
}

func NewRepository(mainDB *pgxpool.Pool, cmsUtil interfaces.CmsUtil, paginator *pagination.Paginator) *repository {
	return &repository{
		mainDB:    mainDB,
		cmsUtil:   cmsUtil,
		paginator: paginator,
	}
}

func (r *repository) GetTableDataTypes(tableName string, columnConverter *func(s string) string) (columns []models.Column, err error) {
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

	if columnConverter != nil {
		converter := *columnConverter
		for i, column := range columns {
			columns[i].ColumnName = converter(column.ColumnName)
		}
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
		fmt.Sprintf(`
		SELECT * FROM %s
		OFFSET %d LIMIT %d;`,
			tableName,
			offset,
			limit,
		),
	)

	if err != nil {
		return
	}

	for i, data := range res {
		res[i] = util.ConvertMapKeys(data, strcase.ToLowerCamel)
	}

	return
}

func (r *repository) GetRelatedTableData(tableName string, relatedTable string) (res []map[string]interface{}, err error) {
	ctx := context.Background()

	res = make([]map[string]interface{}, 0)
	err = pgxscan.Select(
		ctx,
		r.mainDB,
		&res,
		fmt.Sprintf(`
		SELECT %s.* FROM %s
		INNER JOIN %s
		ON %s.id = %s.%s_id;
		`,
			relatedTable,
			tableName,
			relatedTable,
			tableName,
			relatedTable,
			tableName,
		),
	)

	if err != nil {
		return
	}

	for i, data := range res {
		res[i] = util.ConvertMapKeys(data, strcase.ToLowerCamel)
	}

	return
}

func (r *repository) CreateTableData(tableName string, data map[string]interface{}, attributes []string) (res map[string]interface{}, err error) {
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

	psqlAttributes := r.cmsUtil.ConvertAttributesToPsqlAttributes(attributes)
	joinedColumns := strings.Join(columns, ",")
	joinedPlaceholders := strings.Join(placeholders, ",")
	joinedPsqlAttributes := strings.Join(psqlAttributes, ",")
	sql := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s) RETURNING %s;",
		tableName,
		joinedColumns,
		joinedPlaceholders,
		joinedPsqlAttributes,
	)

	row := r.mainDB.QueryRow(
		ctx,
		sql,
		values...,
	)

	res, err = r.cmsUtil.ConvertRowToMap(row, attributes)
	return
}

func (r *repository) UpdateTableDataByPk(tableName string, pk int, data map[string]interface{}, attributes []string) (res map[string]interface{}, err error) {
	ctx := context.Background()

	var placeholders []string
	var values []interface{}
	i := 1

	for key, value := range data {
		values = append(values, value)
		placeholders = append(placeholders, fmt.Sprintf("%s = $%d", key, i))
		i++
	}

	psqlAttributes := r.cmsUtil.ConvertAttributesToPsqlAttributes(attributes)

	joinedPlaceholders := strings.Join(placeholders, ",")
	joinedPsqlAttributes := strings.Join(psqlAttributes, ",")

	sql := fmt.Sprintf(
		`UPDATE %s
		SET %s
		WHERE id = %d
		RETURNING %s;`,
		tableName,
		joinedPlaceholders,
		pk,
		joinedPsqlAttributes,
	)

	row := r.mainDB.QueryRow(
		ctx,
		sql,
		values...,
	)

	res, err = r.cmsUtil.ConvertRowToMap(row, attributes)
	return
}

func (r *repository) DeleteTableDataByPk(tableName string, pk int) (err error) {
	ctx := context.Background()

	sql := fmt.Sprintf(
		`DELETE FROM %s
		WHERE id = %d`,
		tableName,
		pk,
	)

	_, err = r.mainDB.Exec(
		ctx,
		sql,
	)

	return
}
