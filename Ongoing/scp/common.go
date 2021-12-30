// 주요내용
// 1. CommonResponse scp향으로 변경 완료
// 2. logResponse 부분은 아예 폐쇄. >> 정상적인 로그 리스폰스 제한 (response관련 매개변수들이 ncp와 다르기 때문에 일단 주석처리)
// 3. containsInStringList 역할 모름 >> 변경 x

package scp

// 대상 csp의 sdk를 import한다. < 이후 scp 향으로 수정 필요 
import (
	"encoding/json"
	"fmt"
	"log"

	"strings"

	"github.com/jhp0204/Provider_test1/scp-clilib"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
)

// ★추가할 Response 확인 >> SCP 기준 ProjectId, requestId, resourceId로 차이 有  >> 
type CommonResponse struct {
	RequestId     *string `json:"requestId,omitempty"`
	projectId    *string `json:"projectId,omitempty"`
	resourceId *string `json:"resourceId,omitempty"`
}

//미수정
type CommonCode struct {
	Code     *string `json:"code,omitempty"`
	CodeName *string `json:"codeName,omitempty"`
}

//CommonError response error body  >> SCP 기준 Error 시, response Error body 확인 필요
type CommonError struct {
	ReturnCode    string
	ReturnMessage string
}

// Response 내용 확인 필요
//func logErrorResponse(tag string, err error, args interface{}) {
//	param, _ := json.Marshal(args)
//	log.Printf("[ERROR] %s error params=%s, err=%s", tag, param, err)
//}

//func logCommonRequest(tag string, args interface{}) {
//	param, _ := json.Marshal(args)
//	log.Printf("[INFO] %s params=%s", tag, param)
//}

//func logResponse(tag string, args interface{}) {
//	resp, _ := json.Marshal(args)
//	log.Printf("[INFO] %s response=%s", tag, resp)
//}

//func logCommonResponse(tag string, commonResponse *CommonResponse, logs ...string) {
//	result := fmt.Sprintf("RequestID: %s, ReturnCode: %s, ReturnMessage: %s", ncloud.StringValue(commonResponse.RequestId), ncloud.StringValue(commonResponse.ReturnCode), ncloud.StringValue(commonResponse.ReturnMessage))
//	log.Printf("[INFO] %s success response=%s %s", tag, result, strings.Join(logs, " "))
//}

//func isRetryableErr(commResp *CommonResponse, code []string) bool {
//	for _, c := range code {
//		if commResp != nil && commResp.ReturnCode != nil && ncloud.StringValue(commResp.ReturnCode) == c {
//			return true
//		}
//	}
//
//	return false
//}

func containsInStringList(str string, s []string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
