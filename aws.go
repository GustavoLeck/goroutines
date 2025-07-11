package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	clientAws *s3.Client
	awsOnce   sync.Once
)

func connectAws() *s3.Client {
	awsOnce.Do(func() {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion("us-east-1"),
			config.WithRetryMaxAttempts(3),
		)
		if err != nil {
			panic(fmt.Sprintf("Erro ao carregar configuração AWS: %v", err))
		}
		clientAws = s3.NewFromConfig(cfg)

		if err := testAwsCredentials(); err != nil {
			panic(fmt.Sprintf("Erro ao testar credenciais AWS: %v", err))
		}

		fmt.Println("=> Conexão com o AWS S3 estabelecida com sucesso!")
	})
	return clientAws
}

func testAwsCredentials() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Testar listando buckets
	_, err := clientAws.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("falha ao testar credenciais: %w", err)
	}

	fmt.Println("✅ Credenciais AWS validadas com sucesso!")
	return nil
}

func sendDataS3(value GetArquivo) error {
	_, err := clientAws.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String("sync-legado"),
		Key:         aws.String(value.Hash),
		Body:        bytes.NewReader(value.Binario),
		ContentType: aws.String(value.MimeType),
	})
	if err != nil {
		fmt.Println("Erro ao enviar arquivo para o S3:", err)
		return err
	}
	return nil
}
