package conftools

import (
	"strings"
	"net/url"
	"fmt"
)

func ParseHost(host string) (*url.URL, error) {
	if !strings.Contains(host, "//") {
		host = "http://" + host
	}
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	if u.Host == "" {
		return nil, fmt.Errorf("no host specified")
	}
	return u, nil
}

func HideUrlCredentials(u *url.URL) *url.URL {
	// Copy original url, replace credentials, e. g. for logging
	if u.User != nil {
		u2 := new(url.URL)
		*u2 = *u
		u2.User = url.UserPassword("xxx", "xxx")
		return u2
	}
	return u
}