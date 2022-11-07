package baidu

import "fmt"

const (
	URL    = "https://www.baidu.com"
	DOMAIN = "www.baidu.com"
	ACCEPT = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
)

func init() {
	fmt.Printf("baidu init,URL:%v\n", URL)
}

func LoadConf() {
	fmt.Println("load baidu conf")
}
