package entities

type User struct {
	ID       string
	Login    string
	Password string
	Balance  float32
}

type UserBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type UserParameters struct {
	ID       string
	Login    string
	Password string
	Balance  float32
}

func (m *UserParameters) IsValid() bool {
	return m.Login != "" || m.Balance >= 0
}

type UserListPars struct {
	Login      *string
	Balance    *float32
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
