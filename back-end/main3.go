package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type A struct {
	Id int64 `json:"id,string"`
}
func main() {
	a := A{
		Id: 9223372036854775807,
	}
	s1, _ := json.Marshal(a)
	fmt.Println(string(s1))

	var b A
	_ = json.Unmarshal(s1, &b)
	fmt.Println(b, reflect.TypeOf(b.Id))

	//
	const c = `{"id": "9223372036854775807"}` // `{"id": 9223372036854775807}`会导致解析JSON失败！
	var d A
	_ = json.Unmarshal([]byte(c), &d)
	fmt.Println(d.Id, reflect.TypeOf(d.Id))
}
