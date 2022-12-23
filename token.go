package hilib

import (
	"encoding/xml"
	"io"
	"net/http"
)

type ResToken struct {
	Token int `xml:"token"`

	raw string `xml:"-"`
}

func (rs *ResToken) Raw() string {
	return rs.raw
}

func (rs *ResToken) setRaw(str string) {
	rs.raw = str
}

type ReqToken struct {
}

func (rs *ReqToken) ReqPath() string {
	return "api/webserver/token"
}

func (rs *ResToken) ReqPath() string {
	return "api/webserver/token"
}

func (rs *ReqToken) Request(c *Config) (r Response, err error) {
	hr, err := http.Get(c.BaseURL + rs.ReqPath())
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(hr.Body)
	if err != nil {
		return nil, err
	}

	var res ResToken
	err = xml.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	res.setRaw(string(body))

	return &res, nil
}
