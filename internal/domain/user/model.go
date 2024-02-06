package user

type User struct {
	ID       string
	Login    string
	Password string
	Balance  int
}

type UserBalance struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

type GetPars struct {
	ID       string
	Login    string
	Password string
	Balance  int
}

func (m *GetPars) IsValid() bool {
	return m.Login != "" || m.Balance >= 0
}

type ListPars struct {
	Login      *string
	Balance    *int
	MinBalance *int
	MaxBalance *int
}

//type Edit struct {
//	ID       *string
//	Login    *string
//	Password *string
//	Balance  *int
//
//	PrevValue *User
//}
