package RuleClient

import (
	"crypto/tls"
	"github.com/P001water/P1finger/libs/p1httputils"
	"net/http"
	"time"
)

type RuleClient struct {
	DefaultFingerPath string                 // 默认指纹库路径
	P1FingerPrints    FingerPrintsTdSafeType // P1finger指纹库
	RuleMode          string

	ProxyUrl              string       // 可选 代理地址
	ProxyClient           *http.Client // http客户端 默认跟随重定向
	ProxyNoRedirectCilent *http.Client // http客户端 默认禁止跟随重定向

	OutputFormat string // 可选 输出格式

	DetectRstTdSafe DetectResultTdSafeType
}

type RuleClientBuilder struct {
	defaultFingerPath     string
	useDefaultFingerFiles bool
	ruleMode              string
	outputFormat          string
	timeout               time.Duration
	proxyURL              string
}

func NewRuleClientBuilder() *RuleClientBuilder {
	return &RuleClientBuilder{
		ruleMode:     "redteam",
		outputFormat: "csv", // 默认值
		timeout:      5 * time.Second,
	}
}

var FingerFilesMap = map[string][]string{
	"full":    {"Firewall.yaml,", "MailServer.yaml", "lowLevel.yaml", "fofa_fingerprints.yaml", "oaSystem.yaml", "other.yaml", "p1_fingerprints.yaml", "supply.yaml", "webApp.yaml"},
	"redteam": {"Firewall.yaml,", "MailServer.yaml", "oaSystem.yaml", "webApp.yaml"},
}

func (b *RuleClientBuilder) Build() (_ *RuleClient, err error) {
	r := &RuleClient{
		RuleMode:     b.ruleMode,
		OutputFormat: b.outputFormat,
		ProxyUrl:     b.proxyURL,
	}

	err = r.LoadFingersFromEXEFS(FingerFilesMap[r.RuleMode])
	if err != nil {
		return r, err
	}

	r.ProxyClient = p1httputils.NewHttpClientBuilder().
		WithProxy(b.proxyURL).
		Build()

	r.ProxyNoRedirectCilent = p1httputils.NewHttpClientBuilder().
		WithProxy(b.proxyURL).
		NoRedirect().
		Build()

	return r, nil
}

func (r *RuleClient) newProxyClientWithTimeout(timeout time.Duration) {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	transCfg := &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		DisableKeepAlives: true,
	}
	r.ProxyClient = &http.Client{
		Timeout:   timeout,
		Transport: transCfg,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (b *RuleClientBuilder) WithProxyURL(url string) *RuleClientBuilder {
	b.proxyURL = url
	return b
}

func (b *RuleClientBuilder) WithRuleMode(ruleMode string) *RuleClientBuilder {
	b.ruleMode = ruleMode
	return b
}

func (b *RuleClientBuilder) WithOutputFormat(format string) *RuleClientBuilder {
	b.outputFormat = format
	return b
}

func (b *RuleClientBuilder) WithTimeout(t time.Duration) *RuleClientBuilder {
	b.timeout = t
	return b
}
