package baidu

import (
	"fmt"
	"gosearch/module/site"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

const SEARCH = URL + "/s?wd=%s" + "&usm=3&rsv_idx=2&rsv_page=1"

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
}

func S(q string) (r string) {
	return Search(q).body
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

func (req *Req) send() (resp *Resp, err error) {
	r := &Resp{code: 200, body: ""}

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

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
	}
	r.body = string(body)

	return r, nil
}
