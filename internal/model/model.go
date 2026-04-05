package model

type WebResponse[T any] struct {
	Code    int           `json:"code"`
	Status  string        `json:"status"`
	Message string        `json:"message,omitempty"`
	Data    T             `json:"data,omitempty"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Errors  any           `json:"errors,omitempty"`
}

type PaginatedWebResponse[T any] struct {
	Items    []T
	Metadata PageMetadata
}

type PageMetadata struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}
