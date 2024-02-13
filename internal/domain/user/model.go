package user

import "context"

type User struct {
	ID       string
	Login    string
	Password string
	Balance  float32
}

type repoDBI interface {
	List(ctx context.Context, pars *ListPars) ([]*User, error)
	Create(ctx context.Context, obj *User) error
	Get(ctx context.Context, pars *GetPars) (*User, error)
	Update(ctx context.Context, pars *GetPars) error
	Delete(ctx context.Context, pars *GetPars) error
	Exists(ctx context.Context, login string) (bool, error)
}

type UserBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type GetPars struct {
	ID       string
	Login    string
	Password string
	Balance  float32
}

func (m *GetPars) IsValid() bool {
	return m.Login != "" || m.Balance >= 0
}

type ListPars struct {
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
