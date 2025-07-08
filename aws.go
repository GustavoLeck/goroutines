package main

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var clientAws *s3.Client

func connectAws() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic("Erro ao conectar ao AWS: " + err.Error())
	}
	clientAws := s3.NewFromConfig(cfg)
	return clientAws
}

func sendDataS3(value GetArquivo) error {
	_, err := clientAws.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String("sync-legado"),
		Key:         aws.String(value.Hash),
		Body:        bytes.NewReader(value.Binario),
		ContentType: aws.String(value.MimeType),
	})
	if err != nil {
		println("Erro ao enviar arquivo para o S3:", value.Hash)
		return err
	}
	return nil
}
