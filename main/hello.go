package main

import (
	"encoding/json"
	"fmt"
	"gosearch/module/site"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
)

type EntityWrapper struct { //注意此处
	entity []site.Entity
	by     func(p, q *site.Entity) bool
}

func (ew EntityWrapper) Len() int { // 重写 Len() 方法
	return len(ew.entity)
}
func (ew EntityWrapper) Swap(i, j int) { // 重写 Swap() 方法
	ew.entity[i], ew.entity[j] = ew.entity[j], ew.entity[i]
}
func (ew EntityWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return ew.by(&ew.entity[i], &ew.entity[j])
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/health", health)
	http.HandleFunc("/search", search)
	site.LoadConf()
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
	_, _ = io.WriteString(w, "ok")
}

func search(w http.ResponseWriter, request *http.Request) {

	fmt.Println(request.URL)
	_ = request.ParseForm()
	q := request.Form.Get("q")
	q = url.QueryEscape(q)
	fmt.Printf("查询内容:%s\n", q)
	jsonResult := &site.JsonResult{Code: 0, Data: &site.EntityList{
		Index: 0,
		Size:  0,
		List:  []site.Entity{},
	}}

	array := [...]site.SearchEngine{
		&site.Wx{Req: site.Req{Q: q}},
		//&site.Google{Req: site.Req{Q: q}},
		&site.Bing{Req: site.Req{Q: q}},
		//&site.Baidu{Req: site.Req{Q: q}},
	}

	for _, engine := range array {
		result := engine.(site.SearchEngine).Search()

		jsonResult.Data.Size += result.Size
		for i, entity := range result.List {
			//初始化自然排序
			entity.PositionScore = (len(result.List) - i) * site.GetPositionWeight(entity.From)
			entity.SearchScore = site.GetSearchScore(entity.From)
			entity.DomainScore = site.GetDomainScore(entity.Host)
			entity.Score = entity.PositionScore + entity.SearchScore + entity.DomainScore
			jsonResult.Data.List = append(jsonResult.Data.List, entity)
		}
	}

	// sort score
	sort.Sort(EntityWrapper{jsonResult.Data.List, func(p, q *site.Entity) bool {
		// Score 递减排序
		return p.Score > q.Score
	}})

	body, err := json.Marshal(jsonResult)
	if err != nil {
		jsonResult.Code = -1
		jsonResult.Msg = err.Error()
		w.WriteHeader(500)
		v, _ := json.Marshal(jsonResult)
		_, _ = w.Write(v)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}
