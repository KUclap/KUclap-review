package storage

import (

	"log"
	"context"
	"kuclap-review-api/src/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/config"	
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"

)

type S3PresignGetObjectAPI interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

func GetDownloadPresignedURL(c context.Context, api S3PresignGetObjectAPI, input *s3.GetObjectInput) (*v4.PresignedHTTPRequest, error) {
	return api.PresignGetObject(c, input)
}

func PresignedURLDownloadRecapS3(recap *models.ResRecap) (string, error) {

	cfg, err		:=	config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Println("[ERR] aws configuration error", err)
		panic("configuration error, " + err.Error())
	}

	client			:=	s3.NewFromConfig(cfg)
	
	classFolder		:=	recap.ClassID + "/"
	key				:=	RECAP_FOLDER + classFolder + recap.FileName

	input			:=	&s3.GetObjectInput{
		Bucket:		aws.String(BUCKET_NAME),
		Key:		aws.String(key),
	}

	pressignedClient		:=	s3.NewPresignClient(client)
	presigned, err			:=	GetDownloadPresignedURL(context.TODO(), pressignedClient, input)
	
	return presigned.URL, err

}