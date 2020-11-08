package routes

import ( 
	"kuclap-review-api/src/dao"
)

var mgoDAO *dao.SessionDAO

func InjectAdapterDAO(sessionDAO *dao.SessionDAO) {
	mgoDAO = sessionDAO
}