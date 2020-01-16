package main

import (
	"net/http"
	"net/url"
)

// Redirect is a HTTP redirect.
type Redirect struct {
	Status   int
	Location *url.URL
}

// MultipleChoices creates a 300 redirect.
func MultipleChoices(location *url.URL) *Redirect {
	return &Redirect{http.StatusMultipleChoices, location}
}

// MovedPermanently creates a 301 redirect.
func MovedPermanently(location *url.URL) *Redirect {
	return &Redirect{http.StatusMovedPermanently, location}
}

// Found creates a 302 redirect.
func Found(location *url.URL) *Redirect {
	return &Redirect{http.StatusFound, location}
}

// SeeOther creates a 303 redirect.
func SeeOther(location *url.URL) *Redirect {
	return &Redirect{http.StatusSeeOther, location}
}

// NotModified creates a 304 redirect.
func NotModified(location *url.URL) *Redirect {
	return &Redirect{http.StatusNotModified, location}
}

// UseProxy creates a 305 redirect.
func UseProxy(location *url.URL) *Redirect {
	return &Redirect{http.StatusUseProxy, location}
}

// TemporaryRedirect creates a 307 redirect.
func TemporaryRedirect(location *url.URL) *Redirect {
	return &Redirect{http.StatusTemporaryRedirect, location}
}

// PermanentRedirect creates a 308 redirect.
func PermanentRedirect(location *url.URL) *Redirect {
	return &Redirect{http.StatusPermanentRedirect, location}
}
