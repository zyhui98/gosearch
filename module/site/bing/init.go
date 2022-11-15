package bing

import "fmt"

const (
	URL    = "https://cn.bing.com"
	DOMAIN = "cn.bing.com"
	ACCEPT = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
)

func init() {
	fmt.Printf("bing init,URL:%v\n", URL)
}

func LoadConf() {
	fmt.Println("load bing conf")
}
