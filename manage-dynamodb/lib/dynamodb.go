package lib

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Template struct {
	Html    string `json:"html"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

type Dynamodb struct {
	client *dynamodb.Client
}

func (db *Dynamodb) BatchWriteItem(ctx context.Context, tableName string, items []types.WriteRequest) {
	// 항목 길이 확인
	size := len(items)

	// Dynamodb batch process
	var i int = 0
	for i <= (size / 25) {
		// 저장을 위한 항목 목록
		var extracted []types.WriteRequest
		// 항목 나누기 (Dynamodb Batch로 저장할 수 있는 항목의 수가 최대 25개로 제한되기 때문에 항목을 25개로 나누어 작업 수행)
		if i < (size / 25) {
			extracted = items[(i * 25):((i + 1) * 25)]
		} else {
			extracted = items[(i * 25):]
		}

		// 데이터 쓰기
		_, err := db.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: extracted,
			},
		})
		// 에러 처리
		if err != nil {
			log.Fatalf("[SDK CALL ERROR] %v", err)
		}
		i++
	}
}

func CreateInputForGuide(list []Guide) []types.WriteRequest {
	// 저장하기 위한 요청 데이터
	var request []types.WriteRequest

	// 데이터 가공
	for index, item := range list {
		// 문자열 슬라이스 변환
		var sources []types.AttributeValue
		for _, elem := range item.Sources {
			sources = append(sources, &types.AttributeValueMemberS{Value: elem})
		}

		// 목록에 추가
		request = append(request, types.WriteRequest{
			PutRequest: &types.PutRequest{
				Item: map[string]types.AttributeValue{
					"id":           &types.AttributeValueMemberS{Value: strconv.Itoa(index)},
					"category":     &types.AttributeValueMemberS{Value: item.Category},
					"published_at": &types.AttributeValueMemberN{Value: strconv.Itoa(item.PublishedAt)},
					"sources":      &types.AttributeValueMemberL{Value: sources},
					"title":        &types.AttributeValueMemberS{Value: item.Title},
					"url":          &types.AttributeValueMemberS{Value: item.Url},
				},
			},
		})
	}
	// 가공 데이터 반환
	return request
}

func NewDynamoDB(cfg aws.Config) *Dynamodb {
	return &Dynamodb{
		client: dynamodb.NewFromConfig(cfg),
	}
}
