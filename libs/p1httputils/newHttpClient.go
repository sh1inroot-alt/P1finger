package p1httputils

import (
	"context"
	"crypto/tls"
	"golang.org/x/net/proxy"
	"net"
	"net/http"
	"net/url"
	"time"
)

type HttpClientBuilder struct {
	transport      *http.Transport
	timeout        time.Duration
	followRedirect bool
}

func NewHttpClientBuilder() *HttpClientBuilder {
	return &HttpClientBuilder{
		transport: &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		},
		timeout:        5 * time.Second,
		followRedirect: true,
	}
}

func (b *HttpClientBuilder) Build() *http.Client {
	client := &http.Client{
		Timeout:   b.timeout,
		Transport: b.transport,
	}

	if !b.followRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return client
}

func (b *HttpClientBuilder) WithTimeout(timeout time.Duration) *HttpClientBuilder {
	b.timeout = timeout
	return b
}

func (b *HttpClientBuilder) WithProxy(proxyAddr string) *HttpClientBuilder {
	if proxyAddr == "" {
		return b
	}

	proxyURL, err := url.Parse(proxyAddr)
	if err != nil {
		return b
	}

	switch proxyURL.Scheme {
	case "socks5", "socks5h":
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err == nil {
			b.transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.Dial(network, addr)
			}
		}
	default: // assumes HTTP/HTTPS proxy
		b.transport.Proxy = http.ProxyURL(proxyURL)
	}

	return b
}

func (b *HttpClientBuilder) NoRedirect() *HttpClientBuilder {
	b.followRedirect = false
	return b
}
