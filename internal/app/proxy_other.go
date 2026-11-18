//go:build !windows

package app

// setPageProxy does nothing: WebKitGTK reads http_proxy from the environment,
// as the app's requests do, and macOS's WebKit follows the system's proxy
// settings alone.
func setPageProxy(string) {}
