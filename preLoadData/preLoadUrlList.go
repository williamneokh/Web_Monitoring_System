package preLoadData

import "williamGoInAction1/dataStructure"

func LoadData() {
	webLink := []string{
		"https://www.yahoo.com",
		"https://www.google.com",
		"https://www.singnet123.com.sg",
		"https://www.singnet.com.sg",
		"https://www.amazon.com.sg",
		"https://www.facebook.com",
		"https://www.facebook123.com",
	}
	dataStructure.UserWebLinkMap["admin"] = webLink

}
