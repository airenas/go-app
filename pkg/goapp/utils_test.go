package goapp

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHidePass_Mongo(t *testing.T) {
	assert.Equal(t, "mongodb://mongo:27017", HidePass("mongodb://mongo:27017"))
	assert.Equal(t, "mongodb://l:----@mongo:27017", HidePass("mongodb://l:olia@mongo:27017"))
}

func TestValidateHTTPResp(t *testing.T) {
	tests := []struct {
		name       string
		code       int
		body       string
		l          int
		wantErrStr string
	}{
		{name: "200", code: 200, body: "OK", l: 100, wantErrStr: ""},
		{name: "299", code: 299, body: "OK", wantErrStr: ""},
		{name: "400", code: 400, body: "error", l: 100, wantErrStr: "resp code: 400\nerror"},
		{name: "503", code: 503, body: "error", l: 100, wantErrStr: "resp code: 503\nerror"},
		{name: "503 no err", code: 503, body: "error", l: 0, wantErrStr: "resp code: 503"},
		{name: "400 long", code: 400, body: strings.Repeat("error", 50), l: 100, wantErrStr: "resp code: 400\n" +
			strings.Repeat("error", 50)[:100] + "..."},
		{name: "400 long", code: 400, body: strings.Repeat("error", 50)[:100], l: 100, wantErrStr: "resp code: 400\n" +
			strings.Repeat("error", 50)[:100]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tResp := httptest.NewRecorder()
			tResp.Body = bytes.NewBuffer([]byte(tt.body))
			tResp.Code = tt.code
			err := ValidateHTTPResp(tResp.Result(), tt.l)
			if tt.wantErrStr != "" {
				assert.Equal(t, tt.wantErrStr, err.Error())
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func Test_getBodyStr(t *testing.T) {
	tests := []struct {
		name string
		rd   io.Reader
		l    int
		want string
	}{
		{name: "Empty", rd: strings.NewReader(""), l: 10, want: ""},
		{name: "New line", rd: strings.NewReader("a"), l: 10, want: "\na"},
		{name: "Trim", rd: strings.NewReader(strings.Repeat("a", 20)), l: 10, want: "\n" + strings.Repeat("a", 10) + "..."},
		{name: "Full", rd: strings.NewReader(strings.Repeat("a", 20)), l: 20, want: "\n" + strings.Repeat("a", 20)},
		{name: "Long", rd: strings.NewReader(strings.Repeat("a", 20000)), l: 20000, want: "\n" + strings.Repeat("a", 20000)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBodyStr(tt.rd, tt.l); got != tt.want {
				t.Errorf("getBodyStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{name: "Leaves", args: "olia olia", want: "olia olia"},
		{name: "Changes \\n", args: "olia \nolia\n", want: "olia  olia "},
		{name: "Changes \\r", args: "olia \rolia\r", want: "olia  olia "},
		{name: "Changes \\n\\r", args: "olia \r\nolia\r\n", want: "olia   olia  "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sanitize(tt.args); got != tt.want {
				t.Errorf("Sanitize() = %v, want %v", got, tt.want)
			}
		})
	}
}
