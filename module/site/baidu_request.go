package site

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func (baidu *Baidu) Enable() (enable bool) {
	return GetEnable(BaiduDomain)
}

func (baidu *Baidu) Search() (result *EntityList) {
	baidu.Req.url = baidu.urlWrap()
	log.Printf("req.url: %s\n", baidu.Req.url)
	resp := &Resp{}
	resp, _ = baidu.send()
	baidu.resp = *resp
	result = baidu.toEntityList()
	return result
}

func (baidu *Baidu) urlWrap() (url string) {
	return fmt.Sprintf(BaiduSearch, baidu.Req.Q)
}

func (baidu *Baidu) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}

	if baidu.resp.doc != nil {
		// Find the review items
		//log.Printf("Review doc: %s\n", resp.doc.Text())
		baidu.resp.doc.Find("div[srcid]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("h3").Find("a").Text()
			url := s.AttrOr("mu", "")
			tpl := s.AttrOr("tpl", "")
			if tpl != "se_com_default" {
				return
			}
			subTitle := s.Find(".c-gap-top-small").Find("span").Text()
			entity := Entity{From: BaiduFrom}
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

func (baidu *Baidu) send() (resp *Resp, err error) {
	resp = &Resp{code: 200}

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", baidu.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", BaiduDomain)
	request.Header.Add("Cookie", BaiduCookie)
	request.Header.Add("Accept", BaiduAccept)
	//request.Header.Add("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		panic(err)
	}

	//处理返回结果
	response, _ := client.Do(request)
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
