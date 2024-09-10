package rcdevice

//设备
type Device struct {
	Name          string `json:"name"`
	RemoteAddress string `json:"remote_address"`
	Port          string `json:"port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

func NewDevice() {

}

type ChangedConfigRequest struct {
}
