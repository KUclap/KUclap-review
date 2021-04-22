package storage

import (

	"time"
	"strconv"
	"log"
	"context"
	"kuclap-review-api/src/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"	

)

type PresignedResponse struct {
	Url			string		`json:"url"`
	RecapID		string		`json:"recapId"`
	ClassID		string		`json:"classID"`
	Timestamp	int64		`json:"timestamp"`
	FileName	string		`json:"fileName"`
	Author		string		`json:"author"`
	Tag			string		`json:"tag"`
	TypeFile	string		`json:"type"`
}

func PresignedURLUploadRecapS3(class *models.Class, author string, recapID string) (PresignedResponse, error) {

	cfg, err		:=	config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Println("[ERR] aws configuration error", err)
		return PresignedResponse{}, err
	}

	client			:=	s3.NewFromConfig(cfg)
	
	classFolder		:=	class.ClassID + "/"
	timestamp		:=	time.Now().UTC().Add(7 * time.Hour).Unix()
	strTimestamp		:=	strconv.Itoa(int(timestamp))
	fileName		:=	class.NameEN + " by " + author + " - kuclap.com " + "(" + strTimestamp + "-" + recapID + ").pdf"
	key				:=	RECAP_FOLDER + classFolder + fileName
	tag				:=	"classID=" + class.ClassID + "&timestamp=" + strTimestamp + "&recapID=" + recapID + "&type=recap"

	input			:=	&s3.PutObjectInput{
		Bucket:			aws.String(BUCKET_NAME),
		Key:			aws.String(key),
		Tagging:		aws.String(tag),
		ContentLength:	LIMIT_FILE_SIZE,
	}

	pressignedClient	:=	s3.NewPresignClient(client)
	presigned, err		:=	pressignedClient.PresignPutObject(context.TODO(), input)

	if err != nil {
		log.Println("[ERR] aws PresignPutObject error", err)
		return PresignedResponse{}, err
	}

	payload	:=	PresignedResponse{
		Url:		presigned.URL,
		RecapID:	recapID,
		ClassID:	class.ClassID,
		Timestamp:	timestamp,
		FileName:	fileName,
		Author:		author,
		Tag:		tag,	//value for x-amz-tagging's header
		TypeFile:	"recap",
	}

	return payload, err

}