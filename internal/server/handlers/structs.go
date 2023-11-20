package handlers

type (
	LoginV1Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	LoginV1Resp struct {
		Token  string `json:"token"`
		UserID int64  `json:"userID"`
	}

	CreateV1Req struct {
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Phone     string `json:"phone"`
	}

	UpdateV1Req struct {
		Id        int64  `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetV1Resp struct {
		Id        int64  `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}
)
