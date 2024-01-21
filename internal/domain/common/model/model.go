package model

type ListParams struct {
	Cols           []string
	Page           int64
	PageSize       int64
	WithTotalCount bool
	OnlyCount      bool
	SortName       string
	Sort           []string
}
