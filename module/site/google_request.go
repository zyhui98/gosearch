package site

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func (g *Google) Enable() (enable bool) {
	return GetEnable(GoogleDomain)
}

func (g *Google) Search() (result *EntityList) {
	g.Req.url = g.urlWrap()
	log.Printf("req.url: %s\n", g.Req.url)
	resp := &Resp{}
	resp, _ = g.send()
	g.resp = *resp
	result = g.toEntityList()
	return result
}

func (g *Google) urlWrap() (url string) {
	return fmt.Sprintf(GoogleSearch, g.Req.Q)
}

func (g *Google) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}

	if g.resp.doc != nil {
		// Find the review items
		//log.Printf("Review doc: %s\n", g.resp.doc.Text())
		g.resp.doc.Find("div[class=MjjYud]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("div[class=yuRUbf]").Find("h3").Text()
			if title == "" {
				return
			}
			url := s.Find("div[class=yuRUbf]").Find("a").AttrOr("href", "")
			subTitle := s.Find("div[class='Z26q7c UK95Uc']").Find("span").Text()
			entity := Entity{From: GoogleFrom}
			entity.Title = title
			entity.SubTitle = subTitle
			entity.Url = url
			host := strings.ReplaceAll(url, "http://", "")
			host = strings.ReplaceAll(host, "https://", "")
			entity.Host = strings.Split(host, "/")[0]
			entityList.List = append(entityList.List, entity)
		})
	}
	return entityList
}

func (g *Google) send() (resp *Resp, err error) {
	resp = &Resp{code: 200}

	uri, err := url.Parse(ProxyURL)
	if err != nil {
		log.Fatalf("url.Parse: %v", err)
	}

	tr = &http.Transport{
		// 设置代理
		Proxy:        http.ProxyURL(uri),
		MaxIdleConns: 100,
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(time.Second * 6)) //设置发送接受数据超时
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", g.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", GoogleDomain)
	request.Header.Add("Cookie", GoogleCookie)
	request.Header.Add("Accept", GoogleAccept)
	request.Header.Add("authority", "www.google.com")
	//request.Header.Add("accept-encoding", "gzip, deflate, br")
	request.Header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	//request.Header.Add("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		panic(err)
	}

	//处理返回结果
	response, err := client.Do(request)
	if response == nil {
		log.Fatalf("response nil: %v", err)
		return nil, nil
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		resp.code = response.StatusCode
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
		return resp, nil
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.code = 200
	resp.doc = doc

	return resp, nil
}
