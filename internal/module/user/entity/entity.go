package entity

type LoginReq struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}

type LoginRes struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	RoleName    string  `json:"role_name"`
	CompanyName string  `json:"company_name"`
	BranchName  *string `json:"branch_name"`
	Token       string  `json:"token"`
}

type UserResult struct {
	Id          string  `db:"id"`
	Name        string  `db:"name"`
	Email       string  `db:"email"`
	RoleName    string  `db:"role_name"`
	CompanyName string  `db:"company_name"`
	BranchName  *string `db:"branch_name"`
	Password    string  `db:"password"`
}
