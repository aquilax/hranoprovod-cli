package shared

import (
	"time"
)

// APINodeList represents list of nodes
type APINodeList []APINode

// APINode contains single node
type APINode struct {
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Calories     float64   `json:"calories"`
	Fat          float64   `json:"fat"`
	Carbohydrate float64   `json:"carbohydrate"`
	Protein      float64   `json:"protein"`
	Barcode      string    `json:"barcode,omitempty"`
	UserID       string    `json:"-"`
	Created      time.Time `json:"created"`
}

// APIError holds service error
type APIError struct {
	IsError bool   `json:"is_error"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (ae APIError) Error() string {
	return ae.Status + ": " + ae.Message
}
