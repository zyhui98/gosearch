package site

import (
	"encoding/json"
	"testing"
)

func TestBing(t *testing.T) {
	engine := &Bing{Req: Req{Q: "yuanbiguo"}}
	s := engine.Search()
	marshal, _ := json.Marshal(s)
	println(string(marshal))
}
