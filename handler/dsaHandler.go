package handler

import (
	"fmt"
	"net/http"
)

type Node struct {
	Username string
	Time     string
	Url      string
	Status   string
	Next     *Node
}

type Stack struct {
	Top  *Node
	Size int
}

type UrlStack struct {
	Url []Stack
}

var UserUrlRecord = make(map[string][]string)
var newStack = Stack{}

func Push(username string, time string, address string, status string) {

	currentNode := newStack.Top
	newNode := &Node{username, time, address, status, nil}

	if currentNode == nil {
		newStack.Top = newNode
	} else {
		tempNode := newStack.Top
		newStack.Top = newNode
		newNode.Next = tempNode
	}
	newStack.Size++
	for _, url := range UserUrlRecord[username] {
		if address == url {

			return
		}

	}

	UserUrlRecord[username] = append(UserUrlRecord[username], address)

	return

}

func PrintAllData(res http.ResponseWriter, req *http.Request) {
	var data []string
	var UserData [][]string

	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	for _, url := range UserUrlRecord[myUser.Username] {
		fmt.Println(url)

		currentNode := newStack.Top

		for currentNode != nil {
			if currentNode.Url == url {
				fmt.Printf("UserID: %v, TimeStamp: %v, URL: %v, Status: %v\n", currentNode.Username, currentNode.Time, currentNode.Url, currentNode.Status)
				data = []string{currentNode.Time, currentNode.Url, currentNode.Status}
				UserData = append(UserData, data)
				//fmt.Println("Is empty!")

			}
			currentNode = currentNode.Next
		}

	}
	_ = tpl.ExecuteTemplate(res, "allrecordeddata.gohtml", UserData)
	return
}
func PrintLatest(res http.ResponseWriter, req *http.Request) {
	var data []string
	var UserData [][]string

	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	for _, url := range UserUrlRecord[myUser.Username] {
		fmt.Println(url)

		currentNode := newStack.Top

		for currentNode != nil {
			if currentNode.Url == url {
				fmt.Printf("UserID: %v, TimeStamp: %v, URL: %v, Status: %v\n", currentNode.Username, currentNode.Time, currentNode.Url, currentNode.Status)
				data = []string{currentNode.Time, currentNode.Url, currentNode.Status}
				UserData = append(UserData, data)
				//fmt.Println("Is empty!")
				break
			}
			currentNode = currentNode.Next
		}

	}
	_ = tpl.ExecuteTemplate(res, "printlatest.gohtml", UserData)
	return
}

type data struct {
	Url     string
	Percent float64
}

func IndividualUrlPerformance(res http.ResponseWriter, req *http.Request) {

	var UserData []data
	myUser := GetUser(res, req)
	if !AlreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	for _, url := range UserUrlRecord[myUser.Username] {
		var totalCount int
		var upStatus int
		var percent float64
		currentNode := newStack.Top

		for currentNode != nil {
			if currentNode.Url == url {
				totalCount++
				if currentNode.Status == "up" {
					upStatus++
				}
			}
			currentNode = currentNode.Next
		}
		if totalCount != 0 {
			percent = (float64(upStatus) / float64(totalCount)) * 100
		}

		fmt.Println(url, percent)
		newData := data{url, percent}
		UserData = append(UserData, newData)
	}
	_ = tpl.ExecuteTemplate(res, "individualurlperformance.gohtml", UserData)
	return
}
