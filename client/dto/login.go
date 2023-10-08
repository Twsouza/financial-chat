package dto

type LoginReq struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginRes struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
