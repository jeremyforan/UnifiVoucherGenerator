package UnifiVoucherGenerator

import (
	"net/http"
	"time"
)

func NoteTimestamp() string {
	currentTime := time.Now()

	readableTimestamp := currentTime.Format("Monday, January 2, 2006 at 3:04 PM")
	return readableTimestamp
}

func addBasicHeaders(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("DNT", "1") // Do Not Track
}

//todo: may need to be moved
//req.Header.Set("Origin", unifiApiBaseUrl)

func addSecurityHeaders(req *http.Request) {
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"122\", \"Not(A:Brand\";v=\"24\", \"Brave\";v=\"122\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"macOS\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-gpc", "1")
}

func NoteTimeStamp() string {
	now := time.Now()           // Get the current date and time
	layout := "January02152004" // Define the custom layout
	return now.Format(layout)   // Format the current time according to the layout
}
