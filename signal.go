package hilib

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type ResSignalStatus struct {
	raw string `xml:"-"`

	PCI     int   `xml:"pci"`
	SC      int   `xml:"sc"`
	Cell_id int   `xml:"cell_id"`
	RSRQ    U_dB  `xml:"rsrq"`
	RSRP    U_dBm `xml:"rsrp"`
	RSSI    U_dBm `xml:"rssi"`
	SINR    U_dB  `xml:"sinr"`
	RSCP    int   `xml:"rscp"`
	ECIO    int   `xml:"ecio"`
	Mode    int   `xml:"mode"`
}

type U_dBm string

func (d U_dBm) Int() (i int, err error) {
	buf := bytes.NewBufferString(string(d))

	_, err = fmt.Fscanf(buf, "%ddBm", &i)
	return
}

type U_dB string

func (d U_dB) Int() (i int, err error) {
	buf := bytes.NewBufferString(string(d))

	_, err = fmt.Fscanf(buf, "%ddB", &i)
	return
}

type ReqSignalStatus struct {
}

func (rs *ReqSignalStatus) ReqPath() string {
	return "api/device/signal"
}

func (rs *ResSignalStatus) ReqPath() string {
	return "api/device/signal"
}

func (rs *ResSignalStatus) Raw() string {
	return rs.raw
}

func (rs *ResSignalStatus) setRaw(str string) {
	rs.raw = str
}

func (rs *ReqSignalStatus) Request(c *Config) (r Response, err error) {
	hr, err := http.Get(c.BaseURL + rs.ReqPath())
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(hr.Body)
	if err != nil {
		return nil, err
	}

	var res ResSignalStatus
	err = xml.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	res.setRaw(string(body))

	return &res, nil
}
