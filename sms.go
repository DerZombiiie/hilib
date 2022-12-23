package hilib

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	//	"net/url"
	"strings"
)

type ReqSmsList struct {
	Token int `xml:"-"`

	PageIndex int `xml:"PageIndex"`
	ReadCount int `xml:"ReadCount"`
	BoxType   int `xml:"BoxType"`
	SortType  int `xml:"SortType"`
	Ascending int `xml:"Ascending"`

	// can only be 0 or 1 so is a bool
	UnreadPreferred int `xml:"UnreadPreferred"`
}

func (*ReqSmsList) ReqPath() string {
	return "api/sms/sms-list"
}

type ResSmsList struct {
	raw string `xml:"-"`
	Req string `xml:"-"`
}

func (*ResSmsList) ReqPath() string {
	return "api/sms/sms-list"
}

func (rsl *ResSmsList) setRaw(str string) {
	rsl.raw = str
}

func (rsl *ResSmsList) Raw() string {
	return rsl.raw
}

func (rs *ReqSmsList) Request(c *Config) (r Response, err error) {
	// convert to valid (and (byte identical to real) request)
	b, err := xml.Marshal(rs)
	if err != nil {
		return nil, err
	}

	// concat with header to make request full
	// wired [:len(xml.Header)-1] to remove last char (\n)
	// and strings.ReplaceAll to rename the outer container to request
	b = []byte(xml.Header[:len(xml.Header)-1] + strings.ReplaceAll(string(b), "ReqSmsList", "request"))

	client := &http.Client{}
	// create a new POST request
	// c.BaseURL is http://192.168.8.1 and rs.ReqPath returns api/sms/sms-list
	req, err := http.NewRequest("POST", c.BaseURL+rs.ReqPath(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Print("") // to not have the compiler complain TODO: remove

	req.TransferEncoding = []string{"identity"} // force to use one packet and not chunked transfer

	// insert token into headers:
	req.Header = http.Header{
		"__RequestVerificationToken": {fmt.Sprintf("%d", rs.Token)},
		"Content-Type":               {"application/x-www-form-urlencoded; charset=UTF-8"},
	}

	// set the body into a ReaderCloser (Close method is nop)
	req.Body = nopCloser{strings.NewReader(string(b))}

	// actually do the request
	hr, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// read response
	body, err := io.ReadAll(hr.Body)
	if err != nil {
		return nil, err
	}

	// parse response (thoughs error about XML being invalid as html <br> tags is not closing)
	var res ResSmsList
	err = xml.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	res.Req = string(b)
	res.setRaw(string(body))

	return &res, err
}
