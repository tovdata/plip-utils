package main

import (
	"context"
	"fmt"

	"workspace/lib"
)

func main() {
	// Context 생성
	ctx := context.TODO()
	// Config 생성
	cfg := lib.Configuration(ctx)
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
