package site

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func (wx *Wx) Enable() (enable bool) {
	return GetEnable(WxDomain)
}

func (wx *Wx) Search() (result *EntityList) {
	wx.Req.url = wx.urlWrap()
	log.Printf("req.url: %s\n", wx.Req.url)
	resp := &Resp{}
	resp, _ = wx.send()
	wx.resp = *resp
	result = wx.toEntityList()
	return result
}

func (wx *Wx) urlWrap() (url string) {
	return fmt.Sprintf(WxSearch, wx.Req.Q)
}

func (wx *Wx) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}

	if wx.resp.doc != nil {
		// Find the review items
		//log.Printf("Wx Review doc: %s\n", wx.resp.doc.Text())
		wx.resp.doc.Find("div[class='txt-box']").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("h3").Find("a").Text()
			url := s.Find("h3").Find("a").AttrOr("href", "")
			url = WxUrl + url
			subTitle := s.Find("p[class='txt-info']").Text()

			entity := Entity{From: WxFrom}
			entity.Title = title
			entity.SubTitle = subTitle
			entity.Url = url
			host := s.Find("div[class='s-p']").Find("a").Text()
			entity.Host = strings.Split(host, "/")[0]
			entityList.List = append(entityList.List, entity)
		})
	}
	return entityList
}

func (wx *Wx) send() (resp *Resp, err error) {

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", wx.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", UserAgent)
	request.Header.Add("Host", WxDomain)
	request.Header.Add("Cookie", WxCookie)
	request.Header.Add("Accept", WxAccept)

	return SendDo(client, request)

}
