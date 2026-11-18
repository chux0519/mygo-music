package app

import (
	"net/url"
	"os"
	"strings"
)

// proxyNames are the variables Go's requests read the proxy from, in the order
// they're asked for: the scheme's first, then any's.
var proxyNames = [...]string{
	"HTTPS_PROXY", "https_proxy",
	"HTTP_PROXY", "http_proxy",
	"ALL_PROXY", "all_proxy",
}

// proxyEnvAtStart is what those variables held as the app started, which the
// listener's setting replaces, and clearing it puts back.
var proxyEnvAtStart = func() [4]string {
	var was [4]string
	for i, name := range proxyNames[:4] {
		was[i] = os.Getenv(name)
	}
	return was
}()

// applyProxy sends the app's requests, and the pages of its windows, through
// the proxy the listener set. With none, it's the one the environment named as
// the app started, so that a page goes where the app's requests go. With
// neither, the environment stays as it was, and pages follow the system's
// settings alone.
//
// It has to run before the first window with a page is made on Windows, whose
// pages read no variable of their own: see setPageProxy.
func (s settings) applyProxy() {
	proxy := strings.TrimSpace(s.Proxy)
	if proxy == "" {
		for _, was := range proxyEnvAtStart {
			if was = strings.TrimSpace(was); was != "" {
				proxy = was
				break
			}
		}
	}
	if proxy == "" || !validProxy(proxy) {
		proxy = ""
	}
	for i, name := range proxyNames[:4] {
		if proxy == "" {
			restoreEnv(name, proxyEnvAtStart[i])
		} else {
			os.Setenv(name, proxy)
		}
	}
	if proxy != "" {
		setPageProxy(proxy)
	}
}

// validProxy reports whether proxy is one the app's requests and pages can use.
func validProxy(proxy string) bool {
	u, err := url.Parse(proxy)
	if err != nil || u.Host == "" {
		return false
	}
	switch u.Scheme {
	case "http", "https", "socks5", "socks5h":
		return true
	}
	return false
}

// restoreEnv puts a variable back to what it held, or stops it being named.
func restoreEnv(name, was string) {
	if was == "" {
		os.Unsetenv(name)
		return
	}
	os.Setenv(name, was)
}
