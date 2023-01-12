package lib

import (
	"context"
	"log"
	"os"

	// AWS
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	// Util
	"github.com/joho/godotenv"
)

func Configuration(ctx context.Context) aws.Config {
	var err error
	// 환경 변수 불러오기
	err = godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("[CONFIG ERROR] Failed to load the enviroment variables, %v\n", err)
	}

	// 환경 변수에서 AWS Credentials 가져오기
	AWS_ACCESS_KEY := os.Getenv("AWS_ACCESS_KEY_ID")
	AWS_SECRET_KEY := os.Getenv("AWS_SECRET_ACCESS_KEY")
	// Credentials 값 유무 확인
	isCredentials := AWS_ACCESS_KEY == "" || AWS_SECRET_KEY == ""

	// 설정 변수
	var cfg aws.Config
	// Configuration
	if isCredentials {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AWS_ACCESS_KEY, AWS_SECRET_KEY, "")))
	} else { // Configuration (Credentials이 존재하지 않을 경우, 기기 내 기본 Credentials으로 진행)
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-2"))
	}
	// 에러 처리
	if err != nil {
		log.Fatalf("[CONFIG ERROR] Unable to load AWS SDK config, %v", err)
	}
	// 설정 값 반환
	return cfg
}
