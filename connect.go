package hilib

import (
	"encoding/xml"
	"fmt"
	"gopkg.in/resty.v1"
	"strings"
)

//go:generate stringer -type DialAction

type DialAction int

const (
	DialDisconnect DialAction = iota
	DialConnect
)

type ReqDial struct {
	Token int `xml:"-"`

	Action DialAction `xml:"Action"`
}

func (rs *ReqDial) ReqPath() string {
	return "api/dialup/dial" // well.... its not dialup (I think :)
}

type ResDial struct {
	ResString
}

func (rd *ResDial) ReqPath() string {
	return "api/dialup/dial"
}

func (rd *ReqDial) Request(c *Config) (r Response, err error) {
	body, err := xml.Marshal(rd)
	if err != nil {
		return nil, err
	}

	// concat with header to make request full
	// wired [:len(xml.Header)-1] to remove last char (\n)
	// and strings.ReplaceAll to rename the outer container to request
	body = []byte(xml.Header[:len(xml.Header)-1] + strings.ReplaceAll(string(body), "ReqSendSMS", "request"))

	resp, err := resty.R().
		SetHeader("__requestverificationtoken", fmt.Sprintf("%d", rd.Token)).
		SetBody(body).
		Post(c.BaseURL + rd.ReqPath())

	var res ResDial
	err = xml.Unmarshal(resp.Body(), &res.ResString.Response)
	if err != nil {
		res.setRaw(string(resp.Body()))
		return &res, err
	}

	res.setRaw(string(resp.Body()))

	return &res, err
}

// Make sure it satisfies the interface
var _ = Request(&ReqDial{})
var _ = Response(&ResDial{})
