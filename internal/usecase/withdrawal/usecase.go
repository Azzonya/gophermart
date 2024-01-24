package withdrawal

type Usecase struct {
	srv WithdrawalServiceI
}

func New(srv WithdrawalServiceI) *Usecase {
	return &Usecase{
		srv: srv,
	}
}
