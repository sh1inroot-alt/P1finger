package RuleClient

import "sync"

type DetectResult struct {
	Host                string `json:"host"`
	OriginUrl           string `json:"origin_target"`
	OriginUrlStatusCode int    `json:"origin_url_status_code"`
	OriginWebTitle      string `json:"web_title"` //Important

	RedirectUrl           string `json:"redirect_url"`
	RedirectUrlStatusCode int    `json:"redirect_url_status_code"`
	RedirectWebTitle      string `json:"redirect_web_title"`

	ContentLength  string   `json:"content_length"`
	SiteUp         string   `json:"site_up"`
	FingerTag      []string `json:"finger_tag"`     //Important
	LastUpdateTime string   `json:"lastupdatetime"` //Important
}

type DetectResultTdSafeType struct {
	mu sync.Mutex
	t  []DetectResult
}

func (s *DetectResultTdSafeType) AddElement(elem DetectResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.t = append(s.t, elem)
}

func (s *DetectResultTdSafeType) GetElements() []DetectResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.t
}
