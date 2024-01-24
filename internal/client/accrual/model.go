package accrual

type SendReq struct {
	Method  string
	NodeId  string
	DbName  string
	Path    string
	Params  map[string]string
	Timeout string
	ReqObj  any
	RepObj  any
}
