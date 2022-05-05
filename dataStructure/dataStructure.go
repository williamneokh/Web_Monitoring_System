package dataStructure

type User struct {
	Username         string
	Password         []byte
	First            string
	Last             string
	Interval         string
	MonitorIsOn      bool
	MonitorIsRunning bool
}

var UserWebLinkMap = make(map[string][]string)
