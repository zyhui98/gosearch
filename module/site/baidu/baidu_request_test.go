package baidu

import (
	"encoding/json"
	"testing"
)

func TestBaidu(t *testing.T) {

	s := S("yuanbiguo")
	marshal, _ := json.Marshal(s)
	println(string(marshal))
}
