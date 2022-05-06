package main

import (
	"fmt"
	"github.com/williamneokh/WebMonitoringSystem/handler"
	"net/http"
)

func init() {
	handler.Initial() // should not appear after go action 1
}

var portNum = ":8080"

func main() {

	http.HandleFunc("/", handler.Index)
	http.HandleFunc("/signup", handler.Signup)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/logout", handler.Logout)
	http.HandleFunc("/dashboard", handler.Dashboard)
	http.HandleFunc("/viewmonitoredurllist", handler.ViewAllMonitoredUrlList)
	http.HandleFunc("/addnewurl", handler.AddNewUrl)
	http.HandleFunc("/deleteurl", handler.DeleteUrl)
	http.HandleFunc("/startstopmonitoring", handler.StartStopMonitoring)
	http.HandleFunc("/startmonitoring", handler.StartMonitoring)
	http.HandleFunc("/allrecordeddata", handler.PrintAllData)
	http.HandleFunc("/printlatest", handler.PrintLatest)
	http.HandleFunc("/individualurlperformance", handler.IndividualUrlPerformance)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	fmt.Println("Server is running on http://localhost" + portNum)
	_ = http.ListenAndServe(portNum, nil)
}
