package storage

import (
	"kuclap-review-api/src/constant"
)


var (
	BUCKET_NAME			string
	RECAP_FOLDER		string
	LIMIT_FILE_SIZE		int64
)

func init() {
	BUCKET_NAME			=	constant.BUCKET_NAME
	RECAP_FOLDER		=	constant.RECAP_FOLDER
	LIMIT_FILE_SIZE		=	constant.LIMIT_FILE_SIZE
}