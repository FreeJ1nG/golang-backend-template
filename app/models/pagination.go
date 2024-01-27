package models

type PaginationOptions struct {
	ItemPerPage int
	CurrentPage int
}

type PaginationMetadata struct {
	LowerBound int
	UpperBound int
	TotalPages int
}
