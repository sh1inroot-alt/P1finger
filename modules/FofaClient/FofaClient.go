/*
env settings:
- FOFA_CLIENT_URL full fofa connnection string, format: <url>/?email=<email>&key=<key>&version=<v2>
- FOFA_SERVER fofa server
- FOFA_EMAIL fofa account email
- FOFA_KEY fofa account key
*/
package FofaClient

import (
	"context"
	"fmt"
	"github.com/projectdiscovery/gologger"
	"net/http"
	"net/url"
)

const (
	defaultServer     = "https://fofa.icu"
	defaultAPIVersion = "v1"
)

// FofaClient of fofa connection
type FofaClient struct {
	Server     string // can set local server for debugging, format: <scheme>://<host>
	APIVersion string // api version
	Email      string // fofa email
	Key        string // fofa key

	Account    AccountInfo // fofa account info
	DeductMode DeductMode  // Deduct Mode

	httpClient *http.Client    //
	Ctx        context.Context // use to cancel requests

	verbose      bool
	onResults    func(results [][]string) // when fetch results callback
	accountDebug bool                     // 调试账号明文信息

}

// Update merge config from config url
func (c *FofaClient) Update(configURL string) error {
	u, err := url.Parse(configURL)
	if err != nil {
		return err
	}

	c.Server = u.Scheme + "://" + u.Host
	if u.Query().Has("email") {
		c.Email = u.Query().Get("email")
	}

	if u.Query().Has("key") {
		c.Key = u.Query().Get("key")
	}

	if u.Query().Has("version") {
		c.APIVersion = u.Query().Get("version")
	}

	return nil
}

// URL generate fofa connection url string
func (c *FofaClient) URL() string {
	return fmt.Sprintf("%s/?email=%s&key=%s&version=%s", c.Server, c.Email, c.Key, c.APIVersion)
}

// GetContext 获取context，用于中止任务
func (c *FofaClient) GetContext() context.Context {
	return c.Ctx
}

// SetContext 设置context，用于中止任务
func (c *FofaClient) SetContext(ctx context.Context) {
	c.Ctx = ctx
}

type ClientOption func(c *FofaClient) error

// WithURL configURL format: <url>/?email=<email>&key=<key>&version=<v2>&tlsdisabled=false&debuglevel=0
func WithURL(configURL string) ClientOption {
	return func(c *FofaClient) error {
		// merge from config
		if len(configURL) > 0 {
			return c.Update(configURL)
		}
		return nil
	}
}

func WithEmail(email string) ClientOption {
	return func(c *FofaClient) error {
		if len(email) > 0 {
			c.Email = email
		}
		return nil
	}
}

func WithApiKey(apiKey string) ClientOption {
	return func(c *FofaClient) error {
		if len(apiKey) > 0 {
			c.Key = apiKey
		}
		return nil
	}
}

// WithOnResults set on results callback
func WithOnResults(onResults func(results [][]string)) ClientOption {
	return func(c *FofaClient) error {
		c.onResults = onResults
		return nil
	}
}

// WithAccountDebug 是否错误里面显示账号密码原始信息
func WithAccountDebug(v bool) ClientOption {
	return func(c *FofaClient) error {
		c.accountDebug = v
		return nil
	}
}

// WithDebug 是否显示debug日志
func WithDebug(verbose bool) ClientOption {
	return func(c *FofaClient) error {
		c.verbose = verbose
		return nil
	}
}

// NewClient from fofa connection string to config
func NewClient(options ...ClientOption) (c *FofaClient, err error) {
	c = &FofaClient{
		Server:     defaultServer,
		APIVersion: defaultAPIVersion,
		Ctx:        context.Background(),
	}

	for _, opt := range options {
		err = opt(c)
		if err != nil {
			return nil, err
		}
	}

	// fetch one time to make sure network is ok
	c.httpClient = &http.Client{}
	c.Account, err = c.AccountInfo()
	if err != nil {
		gologger.Warning().Msgf("account invalid")
		return c, nil
	}

	if c.Account.Error {
		gologger.Warning().Msgf("auth failed")
		return c, fmt.Errorf("auth failed: '%s', make sure key is valid", c.Account.ErrMsg)
	}

	gologger.Info().Msgf("User: %v ApiKey 认证成功. 用户种类: %v", c.Email, c.Account.Category)
	gologger.Info().Msgf("剩余免费F点: %v, API月度剩余查询次数: %v, API月度剩余返回数量: %v", c.Account.RemainFreePoint, c.Account.RemainApiQuery, c.Account.RemainApiData)
	gologger.Info().Msgf("FofaPoint: %v, FCoin: %v", c.Account.FofaPoint, c.Account.FCoin)
	return c, nil
}
