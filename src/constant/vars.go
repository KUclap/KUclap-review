package constant

import (
	"kuclap-review-api/src/config"
)

const (
	RECAP_FOLDER		=	"recap/"
	LIMIT_FILE_SIZE		=	20480	// 1024 * 20 = 20 MB
)

var (
	serverConfig	config.Config
	configuration	config.Configuration
)

var (
	BUCKET_NAME		string
)


func init() {

	serverConfig.Read()
	configuration	=	serverConfig.GetConfig()

	BUCKET_NAME		=	configuration.BucketName
}