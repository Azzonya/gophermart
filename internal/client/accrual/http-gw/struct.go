package http_gw

type RequestResult struct {
	Order   string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}
