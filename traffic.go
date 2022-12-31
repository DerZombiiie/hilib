package hilib

import (
	"encoding/xml"
	"gopkg.in/resty.v1"
)

type ResTrafficStats struct {
	raw string `xml:"-"`

	CurrentConnectTime  int   `xml:"CurrentConnectTime"`
	CurrentUpload       int64 `xml:"CurrentUpload"`
	CurrentDownload     int64 `xml:"CurrentDownload"`
	CurrentDownloadRate int64 `xml:"CurrentDownloadRate"`
	CurrentUploadRate   int64 `xml:"CurrentUploadRate"`
	TotalUpload         int64 `xml:"TotalUpload"`
	TotalDownload       int64 `xml:"TotalDownload"`
	TotalConnectTime    int   `xml:"TotalConnectTime"`

	Showtraffic int `xml:"showtraffic"`
}

type ReqTrafficStats struct {
}

func (rs *ReqTrafficStats) ReqPath() string {
	return "api/monitoring/traffic-statistics"
}

func (rs *ResTrafficStats) ReqPath() string {
	return "api/monitoring/traffic-statistics"
}

func (rs *ResTrafficStats) Raw() string {
	return rs.raw
}

func (rs *ResTrafficStats) setRaw(str string) {
	rs.raw = str
}

func (rs *ReqTrafficStats) Request(c *Config) (r Response, err error) {
	resp, err := resty.R().
		Get(c.BaseURL + rs.ReqPath())
	if err != nil {
		return nil, err
	}

	var res ResTrafficStats
	err = xml.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	res.setRaw(string(resp.Body()))

	return &res, nil
}
