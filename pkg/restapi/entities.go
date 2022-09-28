package restapi

import (
	"net"
	"strings"
	"time"

	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/protocol"
)

// A set of entities to make the conversion to and from API JSON types
// less annoying.

func appDeviceTemplates() map[string]string {
	return map[string]string{
		"application-collection": "/applications",
		"application-data":       "/applications/{aeui}/data{?limit&since}",
		"application-stream":     "/applications/{aeui}/stream",
		"device-collection":      "/applications/{aeui}/devices",
		"device-data":            "/applications/{aeui}/devices/{deui}/data{?limit&since}",
		"gateways":               "/gateways",
		"gateway-info":           "/gateways/{geui}",
	}
}

// apiApplication is the entity used by the REST API for applications
type apiApplication struct {
	ApplicationEUI string `json:"applicationEUI"`
	eui            protocol.EUI
}

// ApplicationList is the list of applications presented by the REST API
type applicationList struct {
	Applications []apiApplication  `json:"applications"`
	Templates    map[string]string `json:"templates"`
}

// NewApplicationList creates a new application list
func newApplicationList() applicationList {
	return applicationList{
		Applications: make([]apiApplication, 0),
		Templates:    appDeviceTemplates(),
	}
}

// NewAppFromModel creates a new application from a model.Application instance
func newAppFromModel(app model.Application) apiApplication {
	return apiApplication{
		ApplicationEUI: app.AppEUI.String(),
		eui:            app.AppEUI,
	}
}

// ToModel converts the API application into a model.Application entity
func (a *apiApplication) ToModel() model.Application {
	return model.Application{
		AppEUI: a.eui,
	}
}

func (a *apiApplication) equals(other apiApplication) bool {
	return a.ApplicationEUI == other.ApplicationEUI
}

// Types of devices; ABP/OTAA
const (
	deviceTypeABP  string = "ABP"
	deviceTypeOTAA string = "OTAA"
)

// APIDevice is the REST API type used for devices
type apiDevice struct {
	DeviceEUI      string `json:"deviceEUI"`
	DevAddr        string `json:"devAddr"`
	AppKey         string `json:"appKey"`
	AppSKey        string `json:"appSKey"`
	NwkSKey        string `json:"nwkSKey"`
	FCntUp         uint16 `json:"fCntUp"`
	FCntDn         uint16 `json:"fCntDn"`
	RelaxedCounter bool   `json:"relaxedCounter"`
	DeviceType     string `json:"deviceType"`
	KeyWarning     bool   `json:"keyWarning"`
	eui            protocol.EUI
	da             protocol.DevAddr
	akey           protocol.AESKey
	askey          protocol.AESKey
	nskey          protocol.AESKey
}

// NewDeviceFromModel creates an APIDevice instance from a model.Device instance.
func newDeviceFromModel(device *model.Device) apiDevice {
	var state = deviceTypeOTAA
	if device.State == model.PersonalizedDevice {
		state = deviceTypeABP
	}
	return apiDevice{
		DeviceEUI:      device.DeviceEUI.String(),
		eui:            device.DeviceEUI,
		DevAddr:        device.DevAddr.String(),
		da:             device.DevAddr,
		akey:           device.AppKey,
		AppKey:         device.AppKey.String(),
		AppSKey:        device.AppSKey.String(),
		askey:          device.AppSKey,
		NwkSKey:        device.NwkSKey.String(),
		nskey:          device.NwkSKey,
		FCntDn:         device.FCntDn,
		FCntUp:         device.FCntUp,
		RelaxedCounter: device.RelaxedCounter,
		DeviceType:     state,
		KeyWarning:     device.KeyWarning,
	}
}

// ToModel converts the instance into model.Device instance
func (d *apiDevice) ToModel(appEUI protocol.EUI) model.Device {
	var state = model.OverTheAirDevice
	if strings.ToUpper(d.DeviceType) == deviceTypeABP {
		state = model.PersonalizedDevice
	}

	return model.Device{
		DeviceEUI:      d.eui,
		DevAddr:        d.da,
		AppKey:         d.akey,
		AppSKey:        d.askey,
		NwkSKey:        d.nskey,
		AppEUI:         appEUI,
		State:          state,
		FCntDn:         d.FCntDn,
		FCntUp:         d.FCntUp,
		RelaxedCounter: d.RelaxedCounter,
		KeyWarning:     d.KeyWarning,
	}
}

// DeviceList is the list of devices
type deviceList struct {
	Devices   []apiDevice       `json:"devices"`
	Templates map[string]string `json:"templates"`
}

// NewDeviceList creates a new device list
func newDeviceList() deviceList {
	return deviceList{
		Devices:   make([]apiDevice, 0),
		Templates: appDeviceTemplates(),
	}
}

// apiGateway is used to convert to and from JSON
type apiGateway struct {
	GatewayEUI string  `json:"gatewayEUI"`
	IP         string  `json:"ip"`
	StrictIP   bool    `json:"strictIP"`
	Latitude   float32 `json:"latitude"`
	Longitude  float32 `json:"longitude"`
	Altitude   float32 `json:"altitude"`
	eui        protocol.EUI
	ipaddr     net.IP
}

// ToModel converts an APIGateway instance to a model.Gateway
func (g *apiGateway) ToModel() model.Gateway {
	eui, _ := protocol.EUIFromString(g.GatewayEUI)
	return model.Gateway{
		GatewayEUI: eui,
		IP:         net.ParseIP(g.IP),
		StrictIP:   g.StrictIP,
		Latitude:   g.Latitude,
		Longitude:  g.Longitude,
		Altitude:   g.Altitude,
	}
}

// NewGatewayFromModel creates a new APIGateway instance from a model.Gateway instance
func newGatewayFromModel(gateway model.Gateway) apiGateway {
	return apiGateway{
		GatewayEUI: gateway.GatewayEUI.String(),
		IP:         gateway.IP.String(),
		StrictIP:   gateway.StrictIP,
		Latitude:   gateway.Latitude,
		Longitude:  gateway.Longitude,
		Altitude:   gateway.Altitude,
	}
}

// GatewayList is the list of gateways
type gatewayList struct {
	Gateways  []apiGateway      `json:"gateways"`
	Templates map[string]string `json:"templates"`
}

// NewGatewayList returns an unpopulated list of gateways
func newGatewayList() gatewayList {
	return gatewayList{
		Gateways: make([]apiGateway, 0),
		Templates: map[string]string{
			"gateway-list": "/gateways",
			"gateway-info": "/gateways/{geui}",
		},
	}
}

// ToUnixMillis converts a nanosecond timestamp into a millisecond timestamp.
// the general assumption is that time.Nanosecond = 1 (which it is)
func ToUnixMillis(unixNanos int64) int64 {
	return unixNanos / int64(time.Millisecond)
}

// FromUnixMillis converts a millisecond timestamp into nanosecond timestamp. Note
// that this assumes that time.Nanosecond = 1 (which it is)
func FromUnixMillis(unixMillis int64) int64 {
	return unixMillis * int64(time.Millisecond)
}

// apiDownstreamMessage is a message that will be sent to a device. The message
// is very similar to the existing model entity but for consistency's sake
// it will be treated like other entities.
type apiDownstreamMessage struct {
	DeviceEUI   string `json:"deviceEUI"`
	Data        string `json:"data"`
	Port        uint8  `json:"port"`
	Ack         bool   `json:"ack"`
	SentTime    int64  `json:"sentTime"`
	CreatedTime int64  `json:"createdTime"`
	AckTime     int64  `json:"ackTime"`
	State       string `json:"state"`
}

// ToModel converts the end-user message into model.DownstreamMessage
func (m *apiDownstreamMessage) ToModel() (model.DownstreamMessage, error) {
	deviceEUI, err := protocol.EUIFromString(m.DeviceEUI)
	if err != nil {
		return model.DownstreamMessage{}, err
	}
	return model.DownstreamMessage{
		DeviceEUI:   deviceEUI,
		Data:        m.Data,
		Port:        m.Port,
		Ack:         m.Ack,
		SentTime:    m.SentTime,
		CreatedTime: m.CreatedTime,
		AckTime:     m.AckTime,
	}, nil
}

func newDownstreamMessageFromModel(msg model.DownstreamMessage) apiDownstreamMessage {
	var state string
	switch msg.State() {
	case model.UnsentState:
		state = "UNSENT"
	case model.SentState:
		state = "SENT"
	case model.AcknowledgedState:
		state = "ACKNOWLEDGED"
	}
	return apiDownstreamMessage{
		DeviceEUI:   msg.DeviceEUI.String(),
		Data:        msg.Data,
		Port:        msg.Port,
		Ack:         msg.Ack,
		SentTime:    msg.SentTime,
		CreatedTime: msg.CreatedTime,
		AckTime:     msg.AckTime,
		State:       state,
	}
}
