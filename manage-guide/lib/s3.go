package lib

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	
	// AWS
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	uploader *manager.Uploader
}

func (storage *S3) UploadFile(ctx context.Context, info Info) string {
	// 파일 이름
	filename := info.FileName
	// 파일 읽기
	data, err := ioutil.ReadFile(strings.Join([]string{"./data/append", filename}, "/"))
	// 에러 처리
	if err != nil {
		log.Fatalf("[IO ERROR] %v", err)
	}
	
	// 확장자 확인 (PDF)
	extension := filepath.Ext(filename)
	// 확장자에 따른 컨텐츠 유형
	var contentType string
	if strings.ToLower(extension) == ".pdf" {
		contentType = "application/pdf"
	} else if strings.ToLower(extension) == ".docx" || strings.ToLower(extension) == ".doc" {
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	}
	
	// 입력 데이터 가공
	input := &s3.PutObjectInput{
		Bucket: aws.String("plip.kr"),
		Key: aws.String(strings.Join([]string{"guide", info.FileName}, "/")),
		Body: bytes.NewReader(data),
		ContentDisposition: aws.String("attachment"),
		ContentType: aws.String(contentType),
	}
	// 업로드
	result, err := storage.uploader.Upload(ctx, input)
	// 에러 처리
	if err != nil {
		log.Fatalf("[UPLOAD ERROR] %v", err)
	}
	// URL 반환
	return result.Location
}

func NewS3(cfg aws.Config) *S3 {
	// 클라이언트 생성
	client := s3.NewFromConfig(cfg)
	// S3 객체 관리를 위한 구조체 반환
	return &S3 {
		uploader: manager.NewUploader(client),
	}
}