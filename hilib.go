package hilib

import (
	"fmt"
	"io"
)

// Config contains the configuration for requests
type Config struct {
	BaseURL string // e.g. "http://192.168.8.1/"
}

type Response interface {
	ReqPath() string
	setRaw(string)
	Raw() string
}

type Request interface {
	Request(c *Config) (res Response, err error)
	ReqPath() string
}

func NewConfig(BaseURL string) *Config {
	return &Config{
		BaseURL: BaseURL,
	}
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func NopCloser(r io.Reader) io.ReadCloser {
	return nopCloser{r}
}

type Bool uint8

const (
	False Bool = 0
	True  Bool = 1
)

func (b Bool) OnOff() string {
	return b.Map("On", "Off")
}

func (b Bool) TrueFalse() string {
	return b.Map("True", "False")
}

func (b Bool) Map(t, f string) string {
	if b == 0 {
		return f
	} else {
		return t
	}
}

func (b Bool) Bool() bool {
	return b == 1
}

type ResString struct {
	raw string `xml:"-"`

	Response string `xml:"response"`
}

func (rs *ResString) setRaw(str string) {
	rs.raw = str
}

func (rs *ResString) Raw() string {
	return rs.raw
}

var DefaultConfig = &Config{
	BaseURL: "http://192.168.8.1/",
}

func GetToken(cnf *Config) (token int, err error) {
	r, err := (&ReqToken{}).Request(cnf)
	if err != nil {
		return
	}

	res, ok := r.(*ResToken)
	if !ok {
		return 0, fmt.Errorf("Token: Casting error!")
	}

	return res.Token, nil
}
