package hilib

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type ConnStatus int

func (cs ConnStatus) String() string {
	switch cs {
	case 902:
		return "disconnected"
	case 901:
		return "connected"

	}

	return fmt.Sprintf("%d", cs)
}

func (rs *ResStatus) IsConnected() bool {
	return rs.ConnStatus == 901
}

// SignalStrengthPercent returns signal strength in % rounded to 1s
func (rs *ResStatus) SignalStrengthPercent() string {
	return fmt.Sprintf("%d%%", int(float64(rs.SignalIcon)/float64(rs.Maxsignal)*100))
}

type WifiConnStatus string
type SignalStrength int
type SignalIcon int
type CurrentNetworkType int
type CurrentServiceDomain int
type RoamingStatus int
type BatteryStatus int
type BatteryLevel int
type SimlockStatus int

const (
	SimlockUnlocked SimlockStatus = iota
	SimlockLocked
)

type SimStatus int

//go:generate stringer -type SimStatus
const (
	SimMissing SimStatus = 255
	// TODO: figure some stuff out
)

type ResStatus struct {
	raw string `xml:"-"`

	ConnStatus           ConnStatus           `xml:"ConnectionStatus"`
	WifiConnStatus       WifiConnStatus       `xml:"WifiConnectionStatus"`
	SignalStrength       SignalStrength       `xml:"SignalStrength"`
	SignalIcon           SignalIcon           `xml:"SignalIcon"`
	CurrentNetworkType   CurrentNetworkType   `xml:"CurrentNetworkType"`
	CurrentServiceDomain CurrentServiceDomain `xml:"CurrentServiceDomain"`
	RoamingStatus        RoamingStatus        `xml:"RoamingStatus"`
	BatteryStatus        BatteryStatus        `xml:"BatteryStatus"`
	BatteryLevel         BatteryLevel         `xml:"BatteryLevel"`
	SimlockStatus        SimlockStatus        `xml:"simlockStatus"`
	WanIPAddress         string               `xml:"WanIPAddress"`
	WanIPv6Address       string               `xml:"WanIPv6Address"`
	PrimaryDns           string               `xml:"PrimaryDns"`
	SecondaryDns         string               `xml:"SecondaryDns"`
	PrimaryIPv6Dns       string               `xml:"PrimaryIPv6Dns"`
	SecondaryIPv6Dns     string               `xml:"SecondaryIPv6Dns"`
	CurrentWifiUser      string               `xml:"CurrentWifiUser"`
	TotalWifiUser        int                  `xml:"TotalWifiUser"`
	Currenttotalwifiuser int                  `xml:"currenttotalwifiuser"`
	ServiceStatus        int                  `xml:"ServiceStatus"`
	SimStatus            SimStatus            `xml:"SimStatus"`
	WifiStatus           int                  `xml:"WifiStatus"`
	CurrentNetworkTypeEx int                  `xml:"CurrentNetworkTypeEx"`
	Maxsignal            int                  `xml:"maxsignal"`
	Wifiindooronly       int                  `xml:"wifiindooronly"`
	Wififrequence        int                  `xml:"wififrequence"`
	MSISDN               string               `xml:"msisdn"`
	Classify             string               `xml:"classify"`
	Flymode              int                  `xml:"flymode"`
}

func (rs *ResStatus) Raw() string {
	return rs.raw
}

func (rs *ResStatus) setRaw(str string) {
	rs.raw = str
}

type ReqStatus struct {
}

func (rs *ReqStatus) ReqPath() string {
	return "api/monitoring/status"
}

func (rs *ResStatus) ReqPath() string {
	return "api/monitoring/status"
}

func (rs *ReqStatus) Request(c *Config) (r Response, err error) {
	hr, err := http.Get(c.BaseURL + rs.ReqPath())
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(hr.Body)
	if err != nil {
		return nil, err
	}

	var res ResStatus
	err = xml.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	res.setRaw(string(body))

	return &res, nil
}
