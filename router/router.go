package router

import (
	"fmt"
	"livego/backends"
	"livego/backends/models"
	"strings"
)

// SelectRoute :This function retrieves the route based on the app name and return the whole route object.
func SelectRoute(url string) (models.Route, bool) {
	urlArray := strings.Split(url, "/")
	// This condition checks for both /[app]/[key] exist in url
	if len(urlArray) < 2 {
		return models.Route{}, false
	}

	// Populate app and key from urlArray
	app := urlArray[0]
	key := urlArray[1]
	r, exist := backends.DB.Select(app)
	if !exist {
		return models.Route{}, false
	}

	// If the Stream field is specified, match on it as well as the app
	if r.Stream != "" && r.Stream != key {
		return models.Route{}, false
	}
	return r, true
}

// GetRoute :This function retrieves the route based on the app name, and reurn an array of endpoint urls
func GetRoute(url string) ([]string, bool) {
	urlArray := strings.Split(url, "/")
	// This condition checks for both /[app]/[key] exist in url
	if len(urlArray) < 2 {
		return nil, false
	}

	// Populate app and key from urlArray
	app := urlArray[0]
	key := urlArray[1]
	r, exist := backends.DB.Select(app)
	if !exist {
		return nil, false
	}

	// If the Stream field is specified, match on it as well as the app
	if r.Stream != "" && r.Stream != key {
		return nil, false
	}
	// Declare empty urlList Array to populate in loop
	var urlList []string

	// Loop through endpoints and build urlList Array
	for _, e := range r.Endpoints {

		// If CopyKey is true, we will use the incoming streamkey in the endpoint as well
		if e.Enabled {
			if r.CopyKey {
				urlList = append(urlList, fmt.Sprintf("rtmp://%v/%v/%v", e.Host, e.App, key))
			} else {
				urlList = append(urlList, fmt.Sprintf("rtmp://%v/%v/%v", e.Host, e.App, e.Stream))
			}
		}

	}

	// Check the length of urlList. If it is 0, return false
	if len(urlList) < 1 {
		return nil, false
	}
	return urlList, true
}

// CheckAppName : Check to make sure the app exist in backend
func CheckAppName(appname string) bool {
	_, exist := backends.DB.Select(appname)
	if !exist {
		return false
	}
	return true
}
