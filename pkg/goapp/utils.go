package goapp

import "net/url"

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
