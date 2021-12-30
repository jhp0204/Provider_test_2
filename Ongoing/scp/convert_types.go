// 이후 scp sdk 개발 후, 상세내용 수정 필요
package scp

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/jhp0204/Provider_test1/scp-clilib"
)

func validElem(i interface{}) bool {
	return reflect.ValueOf(i).Elem().IsValid()
}

func validField(f reflect.Value) bool {
	return (!f.CanAddr() || f.CanAddr() && !f.IsNil()) && f.IsValid()
}

func StringField(f reflect.Value) *string {
	if f.Kind() == reflect.Ptr && f.Type().String() == "*string" {
		return f.Interface().(*string)
	} else if f.Kind() == reflect.Slice && f.Type().String() == "string" {
		return scp.String(f.Interface().(string))
	}
	return nil
}

// SCP의 CommonResponse에 따라 변경 필요 > 변경 진행
func GetCommonResponse(i interface{}) *CommonResponse {
	if i == nil || !validElem(i) {
		return &CommonResponse{}
	}
	var projectId *string
	var requestId *string
	var resourceId *string

	if f := reflect.ValueOf(i).Elem().FieldByName("requestId"); validField(f) {
		requestId = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("projectId"); validField(f) {
		projectId = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("resourceId"); validField(f) {
		resourceId = StringField(f)
	}
	return &CommonResponse{
		requestId:     requestId,
		projectId:    projectId,
		resourceId: resourceId,
	}
}

//GetCommonErrorBody parse common error message
func GetCommonErrorBody(err error) (*CommonError, error) {
	sa := strings.Split(err.Error(), "Body: ")
	var errMsg string

	if len(sa) != 2 {
		return nil, fmt.Errorf("error body is incorrect: %s", err)
	}

	errMsg = sa[1]

	var m map[string]interface{}
	if err := json.Unmarshal([]byte(errMsg), &m); err != nil {
		return nil, err
	}

	e := m["responseError"].(map[string]interface{})

	// 해당 부분은 SCP의 Error return 확인 필요, 공통 오류 코드 기준으로 보았을 , 너무 많이 나온다. 선별해서 여기에 반영시킬 수 있을지 확인 
	return &CommonError{
		ReturnCode:    e["returnCode"].(string),
		ReturnMessage: e["returnMessage"].(string),
	}, nil
}
/* region.go 주석처리에 따른 관련 함수 주석처리
func GetRegion(i interface{}) *Region {
	if i == nil || !reflect.ValueOf(i).Elem().IsValid() {
		return &Region{}
	}
	var regionNo *string
	var regionCode *string
	var regionName *string
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionNo"); validField(f) {
		regionNo = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionCode"); validField(f) {
		regionCode = StringField(f)
	}
	if f := reflect.ValueOf(i).Elem().FieldByName("RegionName"); validField(f) {
		regionName = StringField(f)
	}

	return &Region{
		RegionNo:   regionNo,
		RegionCode: regionCode,
		RegionName: regionName,
	}
}
*/

//StringPtrOrNil return *string from interface{}
func StringPtrOrNil(v interface{}, ok bool) *string {
	if !ok {
		return nil
	}
	return scp.String(v.(string))
}

//Int32PtrOrNil return *int32 from interface{}
func Int32PtrOrNil(v interface{}, ok bool) *int32 {
	if !ok {
		return nil
	}

	switch i := v.(type) {
	case int:
		return scp.Int32(int32(i))
	case int32:
		return scp.Int32(i)
	case int64:
		return scp.Int32(int32(i))
	default:
		return scp.Int32(i.(int32))
	}
}

//BoolPtrOrNil return *bool from interface{}
func BoolPtrOrNil(v interface{}, ok bool) *bool {
	if !ok {
		return nil
	}
	return scp.Bool(v.(bool))
}

// StringListPtrOrNil Convert from interface to []*string
func StringListPtrOrNil(i interface{}, ok bool) []*string {
	if !ok {
		return nil
	}

	// Handling when not slice type
	if r := reflect.ValueOf(i); r.Kind() != reflect.Slice {
		tmp := []interface{}{r.String()}
		i = tmp
	}

	il := i.([]interface{})
	vs := make([]*string, 0, len(il))
	for _, v := range il {
		switch v.(type) {
		case *string:
			vs = append(vs, v.(*string))
		default:
			// TODO: if the value is "" in list, occur crash error.
			vs = append(vs, scp.String(v.(string)))
		}
	}
	return vs
}

//StringOrEmpty Get string from *pointer
func StringOrEmpty(v *string) string {
	if v != nil {
		return *v
	}

	return ""
}

//StringPtrArrToStringArr Convert []*string to []string
func StringPtrArrToStringArr(ptrArray []*string) []string {
	var arr []string
	for _, v := range ptrArray {
		arr = append(arr, *v)
	}

	return arr
}

//SetStringIfNotNilAndEmpty set value map[key] if *string pointer is not nil and not empty
func SetStringIfNotNilAndEmpty(m map[string]interface{}, k string, v *string) {
	if v != nil && len(*v) > 0 {
		m[k] = *v
	}
}

//ConvertToMap convert interface{} to map[string]interface{}
func ConvertToMap(i interface{}) map[string]interface{} {
	if i == nil {
		return nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	var m map[string]interface{}
	json.Unmarshal(b, &m)

	return m
}

//ConvertToArrayMap convert interface{} to map[string]interface{}
func ConvertToArrayMap(i interface{}) []map[string]interface{} {
	if i == nil {
		return nil
	}

	b, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	var m []map[string]interface{}
	json.Unmarshal(b, &m)

	return m
}
