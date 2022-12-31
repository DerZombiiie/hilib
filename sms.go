package hilib

import (
	"encoding/xml"
	"fmt"
	"gopkg.in/resty.v1"
	"strings"
	"sync"
	"time"
)

type BoxType uint8

const (
	BoxInbox  BoxType = 1
	BoxOutbox         = 2
	BoxDraft          = 3
)

type ReqSMSList struct {
	Token int `xml:"-"`

	PageIndex int     `xml:"PageIndex"`
	ReadCount int     `xml:"ReadCount"`
	BoxType   BoxType `xml:"BoxType"`
	SortType  int     `xml:"SortType"`
	Ascending int     `xml:"Ascending"`

	// can only be 0 or 1 so is a bool
	UnreadPreferred int `xml:"UnreadPreferred"`
}

func (*ReqSMSList) ReqPath() string {
	return "api/sms/sms-list"
}

type SMSSaveType int

const ()

//go:generate stringer -type SMSType
type SMSType int

const (
	NormalSMS SMSType = iota
)

// <?xml version="1.0" encoding="UTF-8"?><request><Index>-1</Index><Phones><Phone>015752296712</Phone></Phones><Sca></Sca><Content>Hallo! </Content><Length>7</Length><Reserved>1</Reserved><Date>2022-12-28 19:19:56</Date></request>
type SMSMessage struct {
	SMStat int `xml:"Smstat"`
	Index  int `xml:"Index"`

	// Can be $name of sender (e.g. your ISP / PayPal...)
	Phone    string      `xml:"Phone"`
	Content  string      `xml:"Content"`
	Date     string      `xml:"Date"`
	Sca      string      `xml:"Sca"` // TomSka?
	SaveType SMSSaveType `xml:"SaveType"`
	Priority int         `xml:"Priority"`
	Type     SMSType     `xml:"Smstype"`
}

type ResSMSList struct {
	Count    int          `xml:"Count"`
	Messages []SMSMessage `xml:"Messages>Message"`

	raw string `xml:"-"`
	Req string `xml:"-"`
}

func (*ResSMSList) ReqPath() string {
	return "api/sms/sms-list"
}

func (rsl *ResSMSList) setRaw(str string) {
	rsl.raw = str
}

func (rsl *ResSMSList) Raw() string {
	return rsl.raw
}

func (rs *ReqSMSList) Request(c *Config) (r Response, err error) {
	if rs.Token == 0 {
		rs.Token, err = GetToken(c)
		if err != nil {
			return
		}
	}
	// convert to valid (and (byte identical to real) request)
	body, err := xml.Marshal(rs)
	if err != nil {
		return nil, err
	}

	// concat with header to make request full
	// wired [:len(xml.Header)-1] to remove last char (\n)
	// and strings.ReplaceAll to rename the outer container to request
	body = []byte(xml.Header[:len(xml.Header)-1] + strings.ReplaceAll(string(body), "ReqSmsList", "request"))

	resp, err := resty.R().
		SetHeader("__requestverificationtoken", fmt.Sprintf("%d", rs.Token)).
		SetBody(body).
		Post("http://192.168.8.1/api/sms/sms-list")

	// parse response (thoughs error about XML being invalid as html <br> tags is not closing)
	var res ResSMSList
	err = xml.Unmarshal(resp.Body(), &res)
	if err != nil {
		res.setRaw(string(resp.Body()))
		return &res, err
	}

	res.Req = string(body)
	res.setRaw(string(resp.Body()))

	return &res, err
}

type ReqSendSMS struct {
	Token int `xml:"-"`

	Index    int      `xml:"Index"`
	Phones   []string `xml:"Phones>Phone"`
	Sca      string
	Content  string
	Length   int
	Reserved int
	Date     string `xml:"Date"`
}

type ResSendSMS struct {
	ResString
}

func (*ResSendSMS) ReqPath() string {
	return "api/sms/send-sms"
}

func (*ReqSendSMS) ReqPath() string {
	return "api/sms/send-sms"
}

func (rss *ResSendSMS) setRaw(str string) {
	rss.raw = str
}

func (rss *ResSendSMS) Raw() string {
	return rss.raw
}

func (rs *ReqSendSMS) Request(c *Config) (r Response, err error) {
	if rs.Index == 0 {
		rs.Index = -1
	}

	if rs.Length <= 0 {
		rs.Length = len(rs.Content)
	}

	if rs.Date == "" {
		rs.Date = time.Now().Format("2006-01-02 15:04:05")
	}

	// convert to valid (and (byte identical to real) request)
	body, err := xml.Marshal(rs)
	if err != nil {
		return nil, err
	}

	// concat with header to make request full
	// wired [:len(xml.Header)-1] to remove last char (\n)
	// and strings.ReplaceAll to rename the outer container to request
	body = []byte(xml.Header[:len(xml.Header)-1] + strings.ReplaceAll(string(body), "ReqSendSMS", "request"))

	resp, err := resty.R().
		SetHeader("__requestverificationtoken", fmt.Sprintf("%d", rs.Token)).
		SetBody(body).
		Post(c.BaseURL + rs.ReqPath())

	// parse response (thoughs error about XML being invalid as html <br> tags is not closing)
	var res ResSendSMS
	err = xml.Unmarshal(resp.Body(), &res.ResString.Response)
	if err != nil {
		res.setRaw(string(resp.Body()))
		return &res, err
	}

	res.setRaw(string(resp.Body()))

	return &res, err
}

type ReqSMSCount struct{}

type ResSMSCount struct {
	raw string `xml:"-"`

	LocalUnread  int
	LocalInbox   int
	LocalOutbox  int
	LocalDraft   int
	LocalDeleted int
	SimUnread    int
	SimInbox     int
	SimOutbox    int
	SimDraft     int
	LocalMax     int
	SimMax       int
	SimUsed      int
	NewMsg       int
}

func (rsc *ResSMSCount) ReqPath() string {
	return "api/sms/sms-count"
}

func (rsc *ReqSMSCount) ReqPath() string {
	return "api/sms/sms-count"
}

func (rsc *ResSMSCount) setRaw(str string) {
	rsc.raw = str
}

func (rsc *ResSMSCount) Raw() string {
	return rsc.raw
}

func (rsc *ReqSMSCount) Request(c *Config) (r Response, err error) {
	resp, err := resty.R().
		Get(c.BaseURL + rsc.ReqPath())

	var res ResSMSCount
	err = xml.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	res.setRaw(string(resp.Body()))

	return &res, nil

}

var _ = Request(&ReqSMSCount{})

//
// SMS-Steam
//

var (
	_SMSChansLastInbox   int
	_SMSChansLastInboxMu sync.RWMutex
	_SMSChans            = make([]chan SMSMessage, 0)
	_SMSChansErr         = make([]chan error, 0)
	_SMSChansMu          sync.RWMutex
	_SMSChansOnce        sync.Once
)

func _SMSChanErr(err error) {
	_SMSChansMu.RLock()
	defer _SMSChansMu.RUnlock()

	for k := range _SMSChansErr {
		_SMSChansErr[k] <- err
	}
}

func _SMSChanMsg(msg SMSMessage) {
	_SMSChansMu.RLock()
	defer _SMSChansMu.RUnlock()

	for k := range _SMSChans {
		_SMSChans[k] <- msg
	}
}

// _SMSListen is going to check for new SMS'es every 125ms and will notify all SMSChans
func _SMSListen() {
	t := time.NewTicker(time.Millisecond * 125)

	for {
		<-t.C

		resp, err := (&ReqSMSCount{}).Request(DefaultConfig)
		if err != nil {
			_SMSChanErr(err)

			continue
		}

		r, ok := resp.(*ResSMSCount)
		if !ok {
			continue
		}

		_SMSChansLastInboxMu.Lock()
		if r.LocalInbox > _SMSChansLastInbox {
			// new SMS ($count new ones):
			count := r.LocalInbox - _SMSChansLastInbox
			_SMSChansLastInbox = r.LocalInbox

			resp, err := (&ReqSMSList{
				PageIndex: r.LocalInbox / count,
				ReadCount: count,
				BoxType:   BoxInbox,

				Ascending: 1,
			}).Request(DefaultConfig)

			if err != nil {
				_SMSChanErr(err)

				continue
			}

			r, ok := resp.(*ResSMSList)
			if !ok {
				continue
			}

			for k := range r.Messages {
				_SMSChanMsg(r.Messages[k])
			}
		}
		_SMSChansLastInboxMu.Unlock()
	}
}

// SMSChan makes sure SMS-Updater is running and then registers a new channel to notify for new SMS Messages
func SMSChan() (<-chan SMSMessage, <-chan error) {
	_SMSChansOnce.Do(func() {
		go _SMSListen()
	})

	msgch := make(chan SMSMessage)
	errch := make(chan error)

	_SMSChansMu.Lock()
	defer _SMSChansMu.Unlock()

	_SMSChans = append(_SMSChans, msgch)
	_SMSChansErr = append(_SMSChansErr, errch)

	return msgch, errch
}

type SMSStream struct {
}
