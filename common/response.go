package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type UrlStruct struct {
	Url         string
	HeaderKey   string
	HeaderValue string
}

func (u *UrlStruct) responseGetter() []byte {
	url := u.Url
	client := http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set(u.HeaderKey, u.HeaderValue)

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	return body
}
