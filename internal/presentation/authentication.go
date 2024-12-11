package presentation

type (
	ReqRegisterUser struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RespRegisterUser struct {
		UserID   string `json:"user_id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
)

type (
	ReqLoginUser struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	RespLoginUser struct {
		AccessToken string `json:"access_token"`
		UserID      string `json:"user_id"`
		FullName    string `json:"full_name"`
		Email       string `json:"email"`
	}
)

type UserAuthData struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Entity      string `json:"entity"`
}
