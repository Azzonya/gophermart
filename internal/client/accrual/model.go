package accrual

type SendReq struct {
	Method  string
	Path    string
	Params  map[string]string
	Timeout string
	ReqObj  any
	RepObj  any
}
