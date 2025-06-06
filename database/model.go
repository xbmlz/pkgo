package database

type Order struct {
	Column string `json:"column" form:"column" binding:"required"`
	Desc   bool   `json:"desc" form:"desc"`
}

type PageQuery struct {
	Page     int     `json:"page" form:"page" binding:"required"`
	PageSize int     `json:"page_size" form:"page_size" binding:"required"`
	Orders   []Order `json:"orders" form:"orders"`
}

type PageResult[T any] struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	Records  []T `json:"records"`
}
