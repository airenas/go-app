package goapp

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// HidePass removes pass from URL
func HidePass(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		Log.Warn().Msg("can't parse url.")
		return ""
	}
	if u.User != nil {
		u.User = url.UserPassword(u.User.Username(), "----")
	}
	return u.String()
}

// Estimate estimates and logs execution duration
// sample: defer goapp.Estimate("function")()
func Estimate(name string) func() {
	start := time.Now()
	return func() {
		Log.Info().Str("time", fmt.Sprintf("%v", time.Since(start))).Str("name", name).Msg("took")
	}
}

// ValidateHTTPResp returns error if code is not in [200, 299]
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
			Log.Warn().Err(err).Send()
		}
		if len(bytes) > l {
			// use runes to make sure we don't crack utf-8
			// and drop last symbol
			rns := []rune(string(bytes))
			if len(rns) > 0 {
				rns = rns[:len(rns)-1]
			}
			return "\n" + string(rns) + "..."
		}
		if len(bytes) > 0 {
			return "\n" + string(bytes)
		}
	}
	return ""
}

// Sanitize replaces new lines in str for logging
func Sanitize(str string) string {
	r := strings.NewReplacer("\n", " ", "\r", " ")
	return r.Replace(str)
}

// IsRetryableCode returns true if status is a retryable HTTP code
func IsRetryableCode(c int) bool {
	return c != http.StatusBadRequest && c != http.StatusUnauthorized && c != http.StatusNotFound && c != http.StatusConflict
}

// InvokeWithBackoff func with backoff
func InvokeWithBackoff[K any](ctx context.Context, f func() (K, bool, error), b backoff.BackOff) (K, error) {
	c := 0
	var err error
	var res K
	var retry bool
	op := func() (K, error) {
		select {
		case <-ctx.Done():
			if err != nil {
				return res, backoff.Permanent(err)
			}
			return res, backoff.Permanent(context.DeadlineExceeded)
		default:
			if c > 0 {
				Log.Info().Int("count", c).Msg("retry")
			}
		}
		c++
		res, retry, err = f()
		if err != nil && !retry {
			Log.Info().Msg("not retryable error")
			return res, backoff.Permanent(err)
		}
		return res, err
	}
	return backoff.RetryWithData(op, b)
}

// IsRetryableErr check if err may be retryable
func IsRetryableErr(err error) bool {
	return errors.Is(err, io.EOF) || errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET) ||
		isTimeout(err)
}

func isTimeout(err error) bool {
	e, ok := err.(net.Error)
	return ok && e.Timeout()
}
