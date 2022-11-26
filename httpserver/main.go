package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("os.Getenv(UNIT_ENV):", os.Getenv("APP_ENV"))
	http.HandleFunc("/health", health)
	//handle定义请求访问该服务器里的/healthz路径，就有下面healthz去处理，healthz一般为健康检查
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//定义handle处理函数，只要该healthz被调用，就会写入ok
func health(w http.ResponseWriter, request *http.Request) {
	log.Println(request.URL)
	_ = request.ParseForm()
	log.Println(request.Form.Get("user"))
	io.WriteString(w, "ok")
}
