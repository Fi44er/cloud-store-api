package auth_utils

import (
	"net/http"
)

// Ключи для контекста
type ctxKey string

const (
	IPKey        ctxKey = "ip"
	UserAgentKey ctxKey = "ua"
)

type HeaderInterceptor struct {
	Transport http.RoundTripper
}

func (i *HeaderInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	if ip, ok := req.Context().Value(IPKey).(string); ok {
		req.Header.Set("X-Forwarded-For", ip)
	}
	if ua, ok := req.Context().Value(UserAgentKey).(string); ok {
		req.Header.Set("User-Agent", ua)
	}

	return i.Transport.RoundTrip(req)
}
