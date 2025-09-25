package resp

import "github.com/puoxiu/discron/common/models"

type (
	RspLogin struct {
		User  *models.User `json:"user"`
		Token string       `json:"token"`
	}
)
