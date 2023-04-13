package routes

import (
	"kuclap-review-api/src/dao"
)

var repository *dao.SessionDAO

func InjectAdapterDAO(sessionDAO *dao.SessionDAO) {
	repository = sessionDAO
}
