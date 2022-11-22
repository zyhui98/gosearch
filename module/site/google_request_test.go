package site

import (
	"encoding/json"
	"testing"
)

func TestGoogle(t *testing.T) {
	engine := &Google{Req: Req{Q: "yuanbiguo"}}
	s := engine.Search()
	marshal, _ := json.Marshal(s)
	println(string(marshal))
}
