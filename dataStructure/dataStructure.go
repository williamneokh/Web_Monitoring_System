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

//Data is a struct that will be use in a slice to parse data to the template.
type Data struct {
	IndexNum int
	Time     string
	Url      string
	Status   string
	Percent  float64
}

var UserWebLinkMap = make(map[string][]string)
