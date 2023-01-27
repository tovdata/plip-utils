package lib

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Template struct {
	Html    string `json:"html"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

type SES struct {
	client *sesv2.Client
}

func (ses *SES) DeleteTemplate(ctx context.Context, name string) {
	// SDK 호출
	_, err := ses.client.DeleteEmailTemplate(ctx, &sesv2.DeleteEmailTemplateInput{TemplateName: aws.String(name)})
	if err != nil {
		log.Fatalf("[SDK CALL ERROR] %v\n", err)
	}
}

func (ses *SES) GetTemplate(ctx context.Context, name string) *types.EmailTemplateContent {
	// SDK 호출
	output, err := ses.client.GetEmailTemplate(ctx, &sesv2.GetEmailTemplateInput{TemplateName: aws.String(name)})
	if err != nil {
		log.Fatalf("[SDK CALL ERROR] %v\n", err)
	}
	// 결과 반환
	return output.TemplateContent
}

func (ses *SES) GetTemplates(ctx context.Context) []string {
	// Paginatior 생성
	paginator := sesv2.NewListEmailTemplatesPaginator(ses.client, &sesv2.ListEmailTemplatesInput{PageSize: aws.Int32(25)})
	// Recover
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("[PROCESS ERROR] %v\n", r)
		}
	}()

	// 목록 객체 생성
	var list []string
	// 데이터 조회
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("[PROCESS ERROR] %v\n", err)
		}
		// 데이터 추출
		for _, elem := range output.TemplatesMetadata {
			list = append(list, *elem.TemplateName)
		}
	}
	// 결과 반환
	return list
}

func (ses *SES) SendEmail(ctx context.Context, templateName string, email string) {
	// SDK 호출
	_, err := ses.client.SendEmail(ctx, &sesv2.SendEmailInput{
		Content: &types.EmailContent{
			Template: &types.Template{
				TemplateData: aws.String("{ \"approvalLink\": \"https://www.naver.com\", \"companyName\": \"회사이름\", \"name\": \"사용자이름\" }"),
				TemplateName: aws.String(templateName),
			},
		},
		Destination: &types.Destination{
			ToAddresses: []string{email},
		},
		FromEmailAddress: aws.String("Plip <contact@plip.kr>"),
	})
	// 에러 처리
	if err != nil {
		log.Fatalf("[PROCESS ERROR] %v\n", err)
	}
}

func (ses *SES) SetTemplate(ctx context.Context, filename string, isCreate bool) {
	// 템플릿 객체
	var template Template
	// 설정 파일 읽기
	data, err := ioutil.ReadFile("./templates/" + filename + ".json")
	if err != nil {
		log.Fatalf("[IO ERROR] %v", err)
	}
	// JSON 변환
	json.Unmarshal(data, &template)

	// 템플릿 내용(HTML) 파일 읽기
	content, err := ioutil.ReadFile("./templates/" + filename + ".html")
	if err != nil {
		log.Fatalf("[IO ERROR] %v", err)
	}
	// 템플릿 내용 설정
	template.Html = string(content)

	// 생성 유무에 따른 처리
	if isCreate {
		// 템플릿 생성
		_, err = ses.client.CreateEmailTemplate(ctx, &sesv2.CreateEmailTemplateInput{
			TemplateName: aws.String(template.Name),
			TemplateContent: &types.EmailTemplateContent{
				Html:    aws.String(template.Html),
				Subject: aws.String(template.Subject),
			},
		})
		// 에러 처리
		if err != nil {
			log.Fatalf("[SDK CALL ERROR] %v\n", err)
		}
	} else {
		_, err = ses.client.UpdateEmailTemplate(ctx, &sesv2.UpdateEmailTemplateInput{
			TemplateName: aws.String(template.Name),
			TemplateContent: &types.EmailTemplateContent{
				Html:    aws.String(template.Html),
				Subject: aws.String(template.Subject),
			},
		})
		// 에러 처리
		if err != nil {
			log.Fatalf("[SDK CALL ERROR] %v\n", err)
		}
	}
}

func NewSES(cfg aws.Config) *SES {
	return &SES{
		client: sesv2.NewFromConfig(cfg),
	}
}
