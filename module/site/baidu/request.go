package baidu

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosearch/module/site"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

const SEARCH = URL + "/s?wd=%s" + "&usm=3&rsv_idx=2&rsv_page=1"
const FROM = "百度"

var tr *http.Transport

func init() {
	tr = &http.Transport{
		MaxIdleConns: 100,
		Dial: func(netw, addr string) (net.Conn, error) {
			conn, err := net.DialTimeout(netw, addr, time.Second*2) //设置建立连接超时
			if err != nil {
				return nil, err
			}
			err = conn.SetDeadline(time.Now().Add(time.Second * 3)) //设置发送接受数据超时
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

type Req struct {
	q         string
	url       string
	userAgent string
	http.Cookie
}

type Resp struct {
	code int
	body string
	doc  *goquery.Document
}

type EntityList struct {
	Index int
	Size  int
	List  []Entity
}

type Entity struct {
	Title    string
	Host     string
	Url      string
	SubTitle string
	From     string
}

func S(q string) (result *EntityList) {
	return Search(q)
}

func Search(q string) (result *EntityList) {
	req := &Req{}
	req.q = q
	req.url = req.urlWrap()
	resp, _ := req.send()
	result = resp.toEntityList()
	return result
}

func (req *Req) urlWrap() (url string) {
	return fmt.Sprintf(SEARCH, req.q)
}

func (resp *Resp) toEntityList() (entityList *EntityList) {
	entityList = &EntityList{Index: 0, Size: 10}
	entityList.List = []Entity{}
	if resp.doc != nil {
		// Find the review items
		//fmt.Printf("Review doc: %s\n", resp.doc.Text())
		resp.doc.Find(".c-container").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the Title
			title := s.Find("h3").Find("a").Text()
			url := s.Find("h3").Find("a").AttrOr("href", "")
			subTitle := s.Find(".c-gap-top-small").Find("span").Text()
			if site.Debug {
				fmt.Printf("Review Title: %s\n", title)
				fmt.Printf("Review Url: %s\n", url)
				fmt.Printf("Review SubTitle: %s\n", subTitle)
			}
			entity := Entity{From: FROM}
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

func (req *Req) send() (resp *Resp, err error) {
	resp = &Resp{code: 200}

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", req.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", site.UserAgent)
	request.Header.Add("Host", DOMAIN)
	request.Header.Add("Cookie", site.BaiduCookie)
	request.Header.Add("Accept", ACCEPT)
	//request.Header.Add("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		panic(err)
	}

	//处理返回结果
	response, _ := client.Do(request)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
		return nil, nil
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
