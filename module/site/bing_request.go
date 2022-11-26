package site

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func (bing *Bing) Enable() (enable bool) {
	return GetEnable(BingDomain)
}

func (bing *Bing) Search() (result *EntityList) {
	bing.Req.url = bing.urlWrap()
	log.Printf("req.url: %s\n", bing.Req.url)
	resp := &Resp{}
	resp, _ = bing.send()
	bing.resp = *resp
	result = bing.toEntityList()
	return result
}

func (bing *Bing) urlWrap() (url string) {
	return fmt.Sprintf(BingSearch, bing.Req.Q)
}

func (bing *Bing) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}

	if bing.resp.doc != nil {
		// Find the review items
		//log.Printf("Review doc: %s\n", resp.doc.Text())
		bing.resp.doc.Find("ol#b_results>li[class=b_algo]").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("div[class=b_title]>h2>a").Text()
			url := s.Find("div[class=b_attribution]>cite").Text()
			subTitle := s.Find("div[class=b_caption]>p").Text()
			entity := Entity{From: BingFrom}
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

func (bing *Bing) send() (resp *Resp, err error) {
	resp = &Resp{code: 200}

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", bing.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", BingDomain)
	request.Header.Add("Cookie", BingCoolkie)
	request.Header.Add("Accept", BingAccept)
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
