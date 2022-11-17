package site

import (
	"encoding/json"
	"testing"
)

func TestBaidu(t *testing.T) {
	engine := &Baidu{Req: Req{Q: "yuanbiguo"}}
	s := engine.Search()
	marshal, _ := json.Marshal(s)
	println(string(marshal))
}
