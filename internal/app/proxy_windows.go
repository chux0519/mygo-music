//go:build windows

package app

import (
	"os"
	"strings"
)

// webviewArgsEnv is what WebView2's loader appends to the arguments of the
// browser of every environment the app makes.
const webviewArgsEnv = "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS"

// setPageProxy makes the pages of the app's windows go through proxy, as the
// app's requests do: WebView2 reads no proxy variable of its own, and follows
// the system's settings alone, so its browser is told on the command line.
//
// The variable is the process's, so this is called before the first window
// with a page is made, and a later call changes none of those already made.
func setPageProxy(proxy string) {
	args := make([]string, 0, 2)
	for _, arg := range strings.Fields(os.Getenv(webviewArgsEnv)) {
		if !strings.HasPrefix(arg, "--proxy-server=") {
			args = append(args, arg)
		}
	}
	args = append(args, "--proxy-server="+proxy)
	os.Setenv(webviewArgsEnv, strings.Join(args, " "))
}
