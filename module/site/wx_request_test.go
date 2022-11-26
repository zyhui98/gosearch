package site

import (
	"encoding/json"
	"testing"
)

func TestWx(t *testing.T) {
	engine := &Wx{Req: Req{Q: "口服"}}
	s := engine.Search()
	marshal, _ := json.Marshal(s)
	println(string(marshal))
}
