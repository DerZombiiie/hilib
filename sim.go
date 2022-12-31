package hilib

import (
	"encoding/xml"
	"gopkg.in/resty.v1"
)

//go:generate stringer -type SimState
type SimState int

const (
	SimPinLocked SimState = 260
	SimPukLocked          = 261
)

//go:generate stringer -type PinOptState
type PinOptState int

const (
	PinAvailable PinOptState = 258 // pin-available also puk available
)

type ResSimStatus struct {
	raw string `xml:"-"`

	SimState    SimState    `xml:"SimState"` // TODO: check if same as SimStatus
	PinOptState PinOptState `xml:"PinOptState"`
	SimPinTimes int         `xml:"SimPinTimes"`
	SimPukTimes int         `xml:"SimPukTimes"`
}

func (rs *ResSimStatus) Raw() string {
	return rs.raw
}

func (rs *ResSimStatus) setRaw(str string) {
	rs.raw = str
}

type ReqSimStatus struct {
}

func (rs *ReqSimStatus) ReqPath() string {
	return "api/pin/status"
}

func (rs *ResSimStatus) ReqPath() string {
	return "api/pin/status"
}

func (rs *ReqSimStatus) Request(c *Config) (r Response, err error) {
	resp, err := resty.R().
		Get(c.BaseURL + rs.ReqPath())
	if err != nil {
		return nil, err
	}

	var res ResSimStatus
	err = xml.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	res.setRaw(string(resp.Body()))

	return &res, nil
}

type OperateSimType int

const (
	// If set pin using puk, CurrentPin = NewPin
	OperateSimSetPuk OperateSimType = 4
)

type ReqSimOperation struct {
	raw string `xml:"-"`

	OperateType OperateSimType `xml:"OperateType"`
	CurrentPin  int            `xml:"CurrentPin"`
	NewPin      int            `xml:"NewPin"`
	PukCode     int            `xml:"PukCode"`
}

func (rsp *ReqSimOperation) ReqPath() string {
	return "api/pin/operate"
}

type ResSimOperation struct {
	Response string `xml:"response"`
}

func (rso *ResSimOperation) ReqPath() string {
	return "api/pin/operate"
}

func (rso *ResSimOperation) IsOK() bool {
	return rso.Response == "OK"
}

///*
///* - simlock: api/pin/simlock
///*

type ReqSimLock struct {
}

func (rsl *ReqSimLock) ReqPath() string {
	return "api/pin/simlock"
}

type ResSimLock struct {
	raw string `xml:"-"`

	SimLockEnable      Bool `xml:"SimLockEnable"`
	SimLockRemainTimes int  `xml:"SimLockRemainTimes"`

	// apperantly can be empty
	PSimLockEnable      string `xml:"pSimLockEnable"`
	PSimLockRemainTimes string `xml:"pSimLockRemainTimes"`
}

func (rsl *ResSimLock) setRaw(str string) {
	rsl.raw = str
}

func (rsl *ResSimLock) Raw() string {
	return rsl.raw
}

func (rsl *ResSimLock) ReqPath() string {
	return "api/pin/simlock"
}

func (rs *ReqSimLock) Request(c *Config) (r Response, err error) {
	resp, err := resty.R().
		Get(c.BaseURL + rs.ReqPath())
	if err != nil {
		return nil, err
	}

	var res ResSimLock
	err = xml.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	res.setRaw(string(resp.Body()))

	return &res, nil
}
