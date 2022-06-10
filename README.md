#Web Monitoring System

Few years ago I was tasked to oversee our company's website.
I was also given the task to restart the server if it went down. Many times, it went down during the weekend when no one is
monitoring it. So we have to pay a 3rd party company to provide us the web monitoring services and also pay for sms notification
if the server ever went down.

Having that in my mind, I decided to use this idea to be my Goschool's project. Not only it served it purpose and also will
be a great tool for my work or personal usage.

In this application, I adopt telegram as the notification channel instead of local sms service to keep the cost further reduce.


Application features:
- Monitor any website uptime around the clock
- Multiple website can be monitored at once
- Custom interval timing for website checking
- History records can be view
- Generate uptime report
- Notification via telegram bot

Data Structure:
This application uses 2 data structures to handle 2 sets of data inputs.

The first array data structure being used as storage for url addresses. It also allows the user to add, delete and view all the data.

Second data structure uses a combination of map and stack pointer base. 

Stack method was chosen because of the LIFO characteristics, in this application being able to view the latest record is important, so it would be easily accessible by using a sequential search algorithm to peek at the top of the stack.

Another advantage being the size of the database can grow as much as it needs.
The complexity of this structure will be Map O(1) + Stack O(n) in worst case.

Concurrency:
Checking every single web address to record its data is done by the go concurrency method. The application use http.get to request
a response from web server, while waiting for the web server to response back to application, go concurrency send another http.get request
for another web url address from the slice. Response received from server will be recorded into stack data structure. The process is repeated
until every single address from the slice is checked. For code reference please check internal/urlMonitor/dataStructure.go line 86 to line 114


