package local

const (
	ssid         = "yourSSID"
	pass         = "yourSSIDPassword"
	serverIP     = "yourServerIP"
	serverPort   = 1053
)

func ServerIP() string {
	return serverIP
}

func ServerpPort() int {
	return serverPort
}

func SSID() string {
	return ssid
}

func Password() string {
	return pass
}