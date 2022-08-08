package baidu

import "fmt"

const (
	URL = "https://www.baidu.com"
)

func init() {
	fmt.Printf("baidu init,URL:%v\n", URL)
}

func LoadConf() {
	fmt.Println("load baidu conf")
}
