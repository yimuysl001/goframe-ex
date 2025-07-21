package einit

import "time"

const (
	prikey    = "yimuyslsjqlzh"
	envkeyPre = "goValidity:"
	path      = "config/key"
)

var (
	localPri Pri
)

type Pri struct {
	BootTime time.Time `json:"bootTime"`
	Key      string    `json:"key"`
	LastTime time.Time `json:"lastTime"`
}

func GetLastTime() time.Time {
	return localPri.LastTime
}
