package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func currentIP(u string) (ip string) {
	c := new(http.Client)
	req := NewRequest(c)
	req.Proxy = u
	url := "http://httpbin.org/get"
	resp, _ := req.Get(url)
	d, _ := resp.Json()
	defer resp.Body.Close()

	return d.Get("origin").MustString()
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("{\"origin\": \"127.0.0.1\"}"))
}

func TestHTTPProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(proxyHandler))
	defer proxy.Close()

	httpProxyURL := proxy.URL
	//assert.Equal(t, currentIP(httpProxyURL) == "127.0.0.1", true)
	fmt.Sprintln(currentIP(httpProxyURL) == "127.0.0.1")
}

func TestHTTPSProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(proxyHandler))
	defer proxy.Close()

	httpsProxyURL := proxy.URL
	//assert.Equal(t, currentIP(httpsProxyURL) == "127.0.0.1", true)
	fmt.Sprintln(currentIP(httpsProxyURL) == "127.0.0.1")
}

func TestSocks5Proxy(t *testing.T) {
}
