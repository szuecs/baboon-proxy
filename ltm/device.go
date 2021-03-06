package ltm

import (
	"fmt"
	"net/url"
	"path"

	"github.com/zalando-techmonkeys/baboon-proxy/backend"
	"github.com/zalando-techmonkeys/baboon-proxy/common"
	"github.com/zalando-techmonkeys/baboon-proxy/config"
	"github.com/zalando-techmonkeys/baboon-proxy/errors"
	"github.com/zalando-techmonkeys/baboon-proxy/util"
)

var (
	conf         = config.LoadConfig()
	ltmPartition = conf.Partition["ltm"]
)

// Devices struct provides information
// about loadbalancer cluster itself
type Devices struct {
	Kind  string `json:"kind"`
	Items []struct {
		Device
	} `json:"items"`
}

// Device struct provides information
// about a specific loadbalancer
type Device struct {
	Kind               string   `json:"kind"`
	Name               string   `json:"name"`
	FullPath           string   `json:"fullPath"`
	Generation         int      `json:"generation"`
	ActiveModules      []string `json:"activeModules"`
	BaseMac            string   `json:"baseMac"`
	Build              string   `json:"build"`
	Cert               string   `json:"cert"`
	ChassisID          string   `json:"chassisId"`
	ChassisType        string   `json:"chassisType"`
	ConfigsyncIP       string   `json:"configsyncIp"`
	Edition            string   `json:"edition"`
	FailoverState      string   `json:"failoverState"`
	HaCapacity         int      `json:"haCapacity"`
	Hostname           string   `json:"hostname"`
	Key                string   `json:"key"`
	ManagementIP       string   `json:"managementIp"`
	MarketingName      string   `json:"marketingName"`
	MirrorIP           string   `json:"mirrorIp"`
	MirrorSecondaryIP  string   `json:"mirrorSecondaryIp"`
	MulticastInterface string   `json:"multicastInterface"`
	MulticastIP        string   `json:"multicastIp"`
	MulticastPort      int      `json:"multicastPort"`
	OptionalModules    []string `json:"optionalModules"`
	PlatformID         string   `json:"platformId"`
	Product            string   `json:"product"`
	SelfDevice         string   `json:"selfDevice"`
	TimeZone           string   `json:"timeZone"`
	Version            string   `json:"version"`
	UnicastAddress     []struct {
		EffectiveIP   string `json:"effectiveIp"`
		EffectivePort int    `json:"effectivePort"`
		IP            string `json:"ip"`
		Port          int    `json:"port"`
	} `json:"unicastAddress"`
}

// ShowLTMDevice returns information
// of loadbalancer devices
func ShowLTMDevice(inputURL string) (*backend.Response, *Devices, *errors.Error) {
	// Declaration LTM Device
	ltmdevice := new(Devices)
	deviceURL := util.ReplaceLTMUritoDeviceURI(inputURL)
	fmt.Println(deviceURL)
	res, err := backend.Request(common.GET, deviceURL, &ltmdevice)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmdevice, nil
}

// ShowLTMDeviceName returns information
// of a specific loadbalancer device
func ShowLTMDeviceName(host, inputURL string, ltmDeviceNames map[string]string) (*backend.Response, *Device, *errors.Error) {
	// Declaration LTM Device Name
	value := ltmDeviceNames[host]
	ltmdevicename := new(Device)
	u := new(url.URL)
	u.Scheme = common.Protocol
	u.Path = path.Join(host, common.DeviceURI, value)
	devicenameURL := u.String()
	res, err := backend.Request(common.GET, devicenameURL, &ltmdevicename)
	if err != nil {
		return nil, nil, err
	}
	return res, ltmdevicename, nil
}

//Loadbalancer checks which loadbalancer device is active
//TODO no error handling, function interacting with backend, possible bugs here
func Loadbalancer(lbpair string, ltmDeviceNames map[string]string) (string, *errors.Error) {
	lb01 := lbpair + "01"
	lb02 := lbpair + "02"
	u := new(url.URL)
	u.Scheme = common.Protocol
	u.Path = path.Join(lb02, common.DeviceURI)
	_, obj, err := ShowLTMDeviceName(lb01, u.Path, ltmDeviceNames)
	if err != nil {
		return "", err
	}
	f5url := new(url.URL)
	var p string
	if obj.FailoverState != "active" {
		p = path.Join(lb02, common.LtmURI)
	} else {
		p = path.Join(lb01, common.LtmURI)
	}
	f5url.Scheme = common.Protocol
	f5url.Path = p
	return f5url.String(), nil
}
