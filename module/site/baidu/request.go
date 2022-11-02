package baidu

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gosearch/module/site"
	"io"
	"log"
	"net"
	"net/http"
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
	body io.ReadCloser
}

type EntityList struct {
	index int
	size  int
	list  []Entity
}

type Entity struct {
	title    string
	url      string
	subTitle string
	from     string
}

func S(q string) (r string) {
	return ""
}

func Search(q string) (r *Resp) {
	req := &Req{}
	req.q = q
	req.url = req.urlWrap()
	send, _ := req.send()

	resp := Resp{0, send.body}
	return &resp
}

func (req *Req) urlWrap() (url string) {
	return fmt.Sprintf(SEARCH, req.q)
}

func (resp *Resp) toEntity() (entityList *EntityList) {
	entityList = &EntityList{index: 0, size: 10}
	if resp.body != nil {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(resp.body)
		if err != nil {
			log.Fatal(err)
		}

		// Find the review items
		doc.Find(".c-container").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the title
			title := s.Find("h3").Find("a").Text()
			url := s.Find("h3").Find("a").AttrOr("href", "")
			subTitle := s.Find("c-gap-top-small").Find("span").Text()
			fmt.Printf("Review %d: %s\n", i, title)
			fmt.Printf("Review %d: %s\n", i, url)
			fmt.Printf("Review %d: %s\n", i, subTitle)
		})

	}
	return entityList
}

func (req *Req) send() (resp *Resp, err error) {
	resp = &Resp{code: 200, body: nil}

	client := &http.Client{
		Transport: tr,
	}
	//提交请求
	request, err := http.NewRequest("GET", req.urlWrap(), nil)
	if err != nil {
		log.Println(err)
	}

	//增加header选项
	request.Header.Add("User-Agent", site.USER_AGENT)
	request.Header.Add("X-Requested-With", "XMLHttpRequest")
	if err != nil {
		panic(err)
	}

	//处理返回结果
	response, _ := client.Do(request)
	defer response.Body.Close()
	if response.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", response.StatusCode, response.Status)
	}

	resp.body = response.Body

	return resp, nil
}
