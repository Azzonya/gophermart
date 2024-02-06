package user

import "context"

type User struct {
	ID       string
	Login    string
	Password string
	Balance  int
}

type repoDBI interface {
	List(ctx context.Context, pars *ListPars) ([]*User, error)
	Create(ctx context.Context, obj *GetPars) error
	Get(ctx context.Context, pars *GetPars) (*User, bool, error)
	Update(ctx context.Context, pars *GetPars) error
	Delete(ctx context.Context, pars *GetPars) error
	Exists(ctx context.Context, login string) (bool, error)
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
