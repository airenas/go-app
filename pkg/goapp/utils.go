package goapp

import (
	"fmt"
	"io"
	"net/http"
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

//ValidateHTTPResp returns error if code is not in [200, 299]
// bodyLen - size of bytes to try read body
func ValidateHTTPResp(resp *http.Response, bodyLen int) error {
	if !(resp.StatusCode >= 200 && resp.StatusCode <= 299) {
		return fmt.Errorf("resp code: %d%s", resp.StatusCode, getBodyStr(resp.Body, bodyLen))
	}
	return nil
}

func getBodyStr(rd io.Reader, l int) string {
	if l > 0 {
		bytes, err := io.ReadAll(io.LimitReader(rd, int64(l+1)))
		if err != nil && err != io.EOF {
			Log.Warn(err)
		}
		if len(bytes) > l {
			return "\n" + string(bytes[:l]) + "..."
		}
		if len(bytes) > 0 {
			return "\n" + string(bytes)
		}
	}
	return ""
}
