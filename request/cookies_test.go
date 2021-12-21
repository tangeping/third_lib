package request

import (
	"fmt"
	"net/http"
	"testing"
)

func TestCookies(t *testing.T) {
	c := new(http.Client)
	req := NewRequest(c)
	req.Cookies = map[string]string{
		"key": "value",
		"a":   "123",
	}
	resp, _ := req.Get("http://httpbin.org/cookies")
	d, _ := resp.Json()
	defer resp.Body.Close()

	v := map[string]interface{}{
		"key": "value",
		"a":   "123",
	}
	fmt.Sprintln(d.Get("cookies").MustMap(), v)
	//assert.Equal(t, d.Get("cookies").MustMap(), v)
}
