package main

import (
	"encoding/json"
	"fmt"
	"gosearch/module/site/baidu"
	"io"
	"log"
	"net/http"
)

type JsonResult struct {
	Code int
	Msg  string
	Data interface{}
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/health", health)
	http.HandleFunc("/search", search)
	baidu.LoadConf()
	go func() {
		log.Println("go in")
		defer func() {
			if err := recover(); err != nil {
				log.Println("go err:", err)
			}
		}()
		log.Println("go out")
	}()
	//handle定义请求访问该服务器里的/health路径，就有下面health去处理，health一般为健康检查
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

//定义handle处理函数，只要该health被调用，就会写入ok
func health(w http.ResponseWriter, request *http.Request) {
	fmt.Println(request.URL)
	_ = request.ParseForm()
	fmt.Println(request.Form.Get("user"))
	io.WriteString(w, "ok")
}

func search(w http.ResponseWriter, request *http.Request) {
	fmt.Println(request.URL)
	_ = request.ParseForm()
	q := request.Form.Get("q")
	fmt.Printf("查询内容:%s\n", q)
	resultBaidu := baidu.S(q)
	resultBaiduJson, err := json.Marshal(resultBaidu)
	fmt.Printf("baidu:%s\n", string(resultBaiduJson))
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	jsonResult := &JsonResult{Code: 0}
	jsonResult.Data = resultBaidu
	body, err := json.Marshal(jsonResult)
	if err != nil {
		jsonResult.Code = -1
		jsonResult.Msg = err.Error()
		w.WriteHeader(500)
		v, _ := json.Marshal(jsonResult)
		_, _ = w.Write(v)
		return
	}
	_, _ = w.Write(body)
}
