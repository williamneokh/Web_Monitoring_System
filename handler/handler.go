package handler

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/williamneokh/WebMonitoringSystem/config"
	"github.com/williamneokh/WebMonitoringSystem/dataStructure"
	"github.com/williamneokh/WebMonitoringSystem/preLoadData"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var tpl *template.Template
var mapUsers = map[string]dataStructure.User{}
var mapSessions = map[string]string{}

func Initial() {

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
	bPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost) //should not appear after go action 1
	mapUsers["admin"] = dataStructure.User{"admin", bPassword, "admin", "admin", "0", false, false}
	preLoadData.LoadData()
}

func Index(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}
	myUser := GetUser(res, req)
	_ = tpl.ExecuteTemplate(res, "index.gohtml", myUser)
}

func Signup(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}
	var myUser dataStructure.User
	// process form submission
	if req.Method == http.MethodPost {
		// get form values
		username := req.FormValue("username")
		password := req.FormValue("password")
		firstname := req.FormValue("firstname")
		lastname := req.FormValue("lastname")
		interval := req.FormValue("interval")
		if username != "" {
			// check if username exist/ taken
			if _, ok := mapUsers[username]; ok {
				_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "User Name already taken!")
				return
			}
			// create session
			id := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}
			http.SetCookie(res, myCookie)
			mapSessions[myCookie.Value] = username

			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				http.Error(res, "Internal server error", http.StatusInternalServerError)
				return
			}

			myUser = dataStructure.User{username, bPassword, firstname, lastname, interval, false, false}
			mapUsers[username] = myUser
		}
		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return

	}
	_ = tpl.ExecuteTemplate(res, "signup.gohtml", myUser)
}

func Login(res http.ResponseWriter, req *http.Request) {
	if AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}

	// process form submission
	if req.Method == http.MethodPost {
		username := req.FormValue("username")
		password := req.FormValue("password")
		// check if user exist with username
		myUser, ok := mapUsers[username]
		if !ok {
			http.Error(res, "Username and/or password do not match", http.StatusUnauthorized)
			return
		}
		// Matching of password entered
		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
		if err != nil {
			_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "Username and/or password do not match")
			return
		}
		// create session
		id := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:  "myCookie",
			Value: id.String(),
		}
		http.SetCookie(res, myCookie)
		mapSessions[myCookie.Value] = username
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}

	_ = tpl.ExecuteTemplate(res, "login.gohtml", nil)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	delete(mapSessions, myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)

	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func Dashboard(res http.ResponseWriter, req *http.Request) {
	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if myUser.MonitorIsOn == true {
		if myUser.Interval != "0" {
			if myUser.MonitorIsRunning == false {
				StartMonitoring(res, req)
			}
		}
	}

	_ = tpl.ExecuteTemplate(res, "dashboard.gohtml", myUser)
}

func ViewAllMonitoredUrlList(res http.ResponseWriter, req *http.Request) {
	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	for i := range dataStructure.UserWebLinkMap {
		if myUser.Username == i {
			data := dataStructure.UserWebLinkMap[i]
			_ = tpl.ExecuteTemplate(res, "viewmonitoredurllist.gohtml", data)
			return
		}
	}
	_ = tpl.ExecuteTemplate(res, "viewmonitoredurllist.gohtml", nil)

}

func AddNewUrl(res http.ResponseWriter, req *http.Request) {
	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {

		username := req.FormValue("username")
		urladdress := req.FormValue("urladdress")

		if urladdress == "" {
			http.Error(res, "No URL address was key", http.StatusUnauthorized)
			return
		}
		_, err := http.Get(urladdress)
		if err != nil {

			http.Error(res, "Your address might be wrong or server is down", http.StatusNotFound)
			return

		} else {
			for _, url := range dataStructure.UserWebLinkMap[username] {
				if urladdress == url {
					_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "Address is duplicated!")
					return
				}
			}

			dataStructure.UserWebLinkMap[username] = append(dataStructure.UserWebLinkMap[username], urladdress)
			http.Redirect(res, req, "/viewmonitoredurllist", http.StatusSeeOther)
			return
		}

	}

	_ = tpl.ExecuteTemplate(res, "addnewurl.gohtml", myUser)

}

func DeleteUrl(res http.ResponseWriter, req *http.Request) {
	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodPost {

		if myUser.MonitorIsOn == true {

			_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "Please stop monitoring first before deleting URL.")
			return
		}

		urladdress := req.FormValue("urladdress")

		for _, url := range dataStructure.UserWebLinkMap[myUser.Username] {
			if urladdress != url {
				dataStructure.UserWebLinkMap["tempArr"] = append(dataStructure.UserWebLinkMap["tempArr"], url)

			}

		}

		dataStructure.UserWebLinkMap[myUser.Username] = dataStructure.UserWebLinkMap["tempArr"]
		dataStructure.UserWebLinkMap["tempArr"] = []string{} /// need to clear the temp data

	}
	for i := range dataStructure.UserWebLinkMap {
		if myUser.Username == i {
			data := dataStructure.UserWebLinkMap[i]
			_ = tpl.ExecuteTemplate(res, "deleteurl.gohtml", data)
			return
		}
	}
	_ = tpl.ExecuteTemplate(res, "deleteurl.gohtml", nil)
}

func StartStopMonitoring(res http.ResponseWriter, req *http.Request) {
	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	if req.Method == http.MethodPost {
		var monitorIsOn bool

		interval := req.FormValue("interval")

		intervalInt, err := strconv.Atoi(interval)
		if err != nil {

			_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "Interval time only allow positive number(s) to be entered.")
			return
		}
		if intervalInt < 0 {
			_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "Interval time only allow positive number(s) to be entered.")
			return
		}

		if _, ok := mapUsers[myUser.Username]; !ok {
			_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "User Name not found")
			return
		}

		if len(dataStructure.UserWebLinkMap[myUser.Username]) < 1 {

			_ = tpl.ExecuteTemplate(res, "errorresponse.gohtml", "No available url to start monitor, Please add URL first.")
			return
		}
		if interval == "0" {
			monitorIsOn = false
		} else {
			monitorIsOn = true
		}
		copyMap := mapUsers[myUser.Username]
		myUser = dataStructure.User{copyMap.Username, copyMap.Password, copyMap.First, copyMap.Last, interval, monitorIsOn, copyMap.MonitorIsRunning}
		mapUsers[myUser.Username] = myUser

		if intervalInt > 0 {
			if myUser.MonitorIsRunning == false {
				doTask := make(chan bool)
				defer close(doTask)
				go func() {
					<-doTask
					StartMonitoring(res, req) // runs the task once
				}()
			}
		}

	}

	_ = tpl.ExecuteTemplate(res, "startstopmonitoring.gohtml", myUser)
}

func StartMonitoring(res http.ResponseWriter, req *http.Request) {
	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	if len(dataStructure.UserWebLinkMap[myUser.Username]) == 0 {
		//fmt.Printf("%v, has no link to monitor", myUser.Username)
		return
	}
	for {
		var wg sync.WaitGroup

		for _, address := range dataStructure.UserWebLinkMap[myUser.Username] {
			wg.Add(1)
			go checkLink(myUser.Username, address, &wg)
		}
		wg.Wait()

		interval, _ := strconv.Atoi(mapUsers[myUser.Username].Interval)
		if interval <= 0 {
			copyMap := mapUsers[myUser.Username]
			myUser = dataStructure.User{copyMap.Username, copyMap.Password, copyMap.First, copyMap.Last, copyMap.Interval, copyMap.MonitorIsOn, false}
			mapUsers[myUser.Username] = myUser

			return
		} else {
			copyMap := mapUsers[myUser.Username]
			myUser = dataStructure.User{copyMap.Username, copyMap.Password, copyMap.First, copyMap.Last, copyMap.Interval, copyMap.MonitorIsOn, true}
			mapUsers[myUser.Username] = myUser
			time.Sleep(time.Duration(interval) * time.Minute)
		}

	}
}
func checkLink(username string, address string, group *sync.WaitGroup) {
	defer group.Done()
	_, err := http.Get(address)

	if err != nil {
		fmt.Printf(" %v, %v, down %v\n", time.Now().Format(time.RFC822Z), address, username)
		Push(username, time.Now().Format(time.RFC822Z), address, "down")
		downTime := time.Now().Format(time.Stamp)
		_, _ = http.PostForm("https://api.telegram.org/bot"+config.Token+"/sendMessage?chat_id="+config.GroupID+"&text=Status: Down! - "+address+" Timestamp: "+downTime, nil)
	} else {
		fmt.Printf(" %v, %v, up %v\n", time.Now().Format(time.RFC822Z), address, username)
		Push(username, time.Now().Format(time.RFC822Z), address, "up")

	}
}
