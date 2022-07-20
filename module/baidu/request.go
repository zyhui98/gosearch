package baidu

import (
	"fmt"
	"net/http"
)

const SEARCH = URL + "s?wd=%s"

type Request interface {
	urlWrap() (s string)
}

type HttpRequest struct {
	q         string
	userAgent string
	http.Cookie
}

func (h *HttpRequest) urlWrap() (s string) {
	return fmt.Sprintf(SEARCH, h.q)
}
