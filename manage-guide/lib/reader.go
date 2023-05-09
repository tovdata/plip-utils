package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"time"

	// Xlsx
	"github.com/thedatashed/xlsxreader"
)

type Info struct {
	Category    string   `json:"category"`
	FileName		string	 `json:"filename"`
	PublishedAt string   `json:"published_at"`
	Sources     []string `json:"sources"`
	Title       string   `json:"title"`
}
type Guide struct {
	Category    string   `json:"category"`
	PublishedAt int      `json:"published_at"`
	Sources     []string `json:"sources"`
	Title       string   `json:"title"`
	Url         string   `json:"url"`
}

func ReadGuides(filePath string) []Info {
	// 파일 열기
	data, err := ioutil.ReadFile(filePath)
	// 에러 처리
	if err != nil {
		log.Fatalf("[IO ERROR] %v", err)
	}

	// 결과 추출 데이터
	var result []Info
	// 데이터 변환
	err = json.Unmarshal(data, &result)
	// 에러 처리
	if err != nil {
		log.Fatalf("[TRANSFER ERROR] %v", err)
	}
	// 추출 값 반환
	return result
}

func ReadGuideConfigFile(filePath string) []Guide {
	// 엑셀 파일을 읽기 위한 인스턴스 생성
	xl, err := xlsxreader.OpenFile(filePath)
	// 에러 처리
	if err != nil {
		log.Fatalf("[IO ERROR] %v", err)
	}
	// 함수 종료 시, IO 종료
	defer xl.Close()

	// 데이터 추출
	var list []Guide
	for row := range xl.ReadRows("20230110") {
		// 셀 데이터 존재 유무 확인
		if row.Cells == nil || len(row.Cells) != 5 {
			continue
		}

		// 날짜 데이터 변환
		date, err := time.Parse("2006-01-02", row.Cells[2].Value)
		// 에러 처리
		if err != nil {
			log.Fatalf("[Reader Error] %v", err)
		}

		// 목록에 추출 데이터 추가
		list = append(list, Guide{
			Category:    row.Cells[0].Value,
			PublishedAt: int(date.Unix()),
			Sources:     strings.Split(row.Cells[3].Value, ", "),
			Title:       row.Cells[1].Value,
			Url:         row.Cells[4].Value,
		})
	}
	// 첫 번째 행 제거 및 반환
	if len(list) > 0 {
		return list[1:]
	} else {
		return list
	}
}
