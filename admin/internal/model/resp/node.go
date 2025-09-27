package resp

import "github.com/puoxiu/discron/common/models"

type (
	RspNodeSearch struct {
		models.Node
		JobCount int `json:"job_count"`
	}
)