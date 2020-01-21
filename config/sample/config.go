package sample

const (
	enabled      = true
	ssid         = "yourSSIDhere"
	pass         = "yourpasswordhere"
	mqttServer   = "tcp://your-mqtt-server.tld:1883"
	mqttUser     = ""
	mqttPassword = ""
	deviceName   = "DemoLEDTrack"
	trackChannel = "demo-track"
)

func TelemetryEnabled() bool {
	return enabled
}

func SSID() string {
	return ssid
}

func Password() string {
	return pass
}

func MQTTServer() string {
	return mqttServer
}

func MQTTUser() string {
	return mqttUser
}

func MQTTPassword() string {
	return mqttPassword
}

func DeviceName() string {
	return deviceName
}

func TrackChannel() string {
	return trackChannel
}
