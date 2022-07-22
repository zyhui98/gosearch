package baidu

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const SEARCH = URL + "/s?wd=%s"

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
	get, err := http.Get(req.urlWrap())
	if err != nil {
		log.Println(err)
	}
	defer get.Body.Close()
	body, err := ioutil.ReadAll(get.Body)
	if err != nil {
		log.Println(err)
	}
	r.body = string(body)

	return r, nil
}
