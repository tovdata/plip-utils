package main

import (
	"context"
	"flag"
	"fmt"

	"workspace/lib"
)

func main() {
	// Context 생성
	ctx := context.TODO()

	// 플래그
	typePtr := flag.String("type", "get", "AWS SES 템플릿 관리 명령어\nex) create, delete, get, list, update")
	// 플래그 분석
	flag.Parse()
	// 플래그 확인
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	// Config 생성
	cfg := lib.Configuration(ctx)
	// SES 객체 생성
	ses := lib.NewSES(cfg)

	// 메인 프로세스
	if *typePtr == "create" {
		// 설정 파일 이름
		name := InputValue(*typePtr)

		// 템플릿 생성
		ses.SetTemplate(ctx, name, true)
	} else if *typePtr == "delete" {
		// 삭제하려는 템플릿 이름 입력
		name := InputValue(*typePtr)

		// 템플릿 삭제
		ses.DeleteTemplate(ctx, name)
		// 삭제 알림
		fmt.Println("템플릿 삭제 완료")
	} else if *typePtr == "get" {
		// 조회하려는 템플릿 이름 입력
		name := InputValue(*typePtr)

		// 템플릿 조회
		template := ses.GetTemplate(ctx, name)
		// 조회 결과
		fmt.Println("템플릿 제목: ", *template.Subject)
		fmt.Println("템플릿 HTML 내용: ", *template.Html)
	} else if *typePtr == "list" {
		// 템플릿 조회
		templates := ses.GetTemplates(ctx)
		// 조회 결과
		if len(templates) == 0 {
			fmt.Println("생성된 템플릿이 없습니다.")
		} else {
			fmt.Println("=--= 조회 결과")
			for _, template := range templates {
				fmt.Println(template)
			}
		}
	} else if *typePtr == "update" {
		// 설정 파일 이름
		name := InputValue(*typePtr)

		// 템플릿 갱신
		ses.SetTemplate(ctx, name, false)
	} else {
		flag.Usage()
	}
}

func InputValue(iType string) string {
	// 입력을 위한 메시지
	var message string
	// 유형에 따른 메시지 설정
	switch iType {
	case "create":
		message = "설정 파일 이름: "
	case "delete":
		message = "삭제하려는 템플릿 이름: "
	case "get":
		message = "조회하려는 템플릿 이름: "
	case "update":
		message = "설정 파일 이름: "
	}
	// 입력
	if message != "" {
		var input string = ""
		fmt.Print(message)
		fmt.Scanf("%s", &input)
		// 입력 값 반환
		return input
	} else {
		return ""
	}
}
