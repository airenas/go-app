package goapp

import (
	"net/url"
	"time"
)

//HidePass removes pass from URL
func HidePass(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		Log.Warn("Can't parse url.")
		return ""
	}
	if u.User != nil {
		u.User = url.UserPassword(u.User.Username(), "----")
	}
	return u.String()
}

//Estimate estimates and logs execution duration
// sample: defer goapp.Estimate("function")()
func Estimate(name string) func() {
	start := time.Now()
	return func() {
		Log.Infof("%s took %v", name, time.Since(start))
	}
}
