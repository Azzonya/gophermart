package handler

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *RegisterRequest) Decode() {

}
