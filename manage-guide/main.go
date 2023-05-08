package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"workspace/lib"
)

func main() {
	// Context 생성
	ctx := context.TODO()
	// Config 생성
	cfg := lib.Configuration(ctx)

	// S3 객체 생성
	s3 := lib.NewS3(cfg)
	// Dynamodb 객체 생성
	dynamodb := lib.NewDynamoDB(cfg)

	// 추가하고자 하는 데이터 읽기
	data := lib.ReadGuides("./data/append.json")

	// 가이드 추가 작업 진행
	for _, elem := range data {
		// 날짜 데이터 변환
		date, err := time.Parse("2006-01-02", elem.PublishedAt)
		// 에러 처리
		if err != nil {
			log.Fatalf("[Reader Error] %v", err)
		}

		// 파일 업로드
		url := s3.UploadFile(ctx, elem)
		fmt.Println("업로드 완료")
		// 가이드 생성
		guide := &lib.Guide{
			Category: elem.Category,
			PublishedAt: int(date.Unix()),
			Sources: elem.Sources,
			Title: elem.Title,
			Url: url,
		}
		
		// DB에 저장
		dynamodb.WriteItem(ctx, guide)
	}
}

func Initialization(ctx context.Context) {
	// Dynamodb 객체 생성
	dynamodb := lib.NewDynamoDB(cfg)

	// 가이드 설정 데이터 불러오기
	data := lib.ReadGuideConfigFile("./data/data.xlsx")
	// 입력 데이터 가공
	input := lib.CreateInputForGuide(data)
	// 데이터 쓰기
	dynamodb.BatchWriteItem(ctx, "plip-guide", input)
	// 알림
	fmt.Println("데이터 쓰기 완료")
}