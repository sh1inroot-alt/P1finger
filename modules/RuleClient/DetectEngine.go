package RuleClient

import (
	"regexp"
	"strings"

	"github.com/P001water/P1finger/libs/p1httputils"
)

const (
	Up   = "up"
	Down = "down"
)

func (r *RuleClient) Detect(target string) (DetectRst DetectResult, err error) {

	var P1fingerResps []p1httputils.HttpRespUtils
	DetectRst = DetectResult{}
	DetectRst.Host = target

	fixedUrl, err := CheckHttpPrefix(target)
	if err != nil {
		DetectRst.OriginUrl = target
		DetectRst.OriginWebTitle = "not web"
		DetectRst.SiteUp = Down
		DetectRst.FingerTag = []string{"unknown proto http/https, try manually"}
		r.DetectRstTdSafe.AddElement(DetectRst)
		return
	}

	// 首次访问禁止重定向，手动解析重定向
	resp, respUtils, err := p1httputils.HttpGet(fixedUrl, r.ProxyNoRedirectCilent)
	if err != nil {
		DetectRst.OriginUrl = fixedUrl
		DetectRst.OriginUrlStatusCode = respUtils.StatusCode
		DetectRst.ContentLength = respUtils.ContentLength
		DetectRst.OriginWebTitle = respUtils.WebTitle
		DetectRst.SiteUp = Down
		DetectRst.FingerTag = []string{"WebSite down"}
		r.DetectRstTdSafe.AddElement(DetectRst)
		return
	}
	P1fingerResps = append(P1fingerResps, respUtils)
	DetectRst.OriginUrl = respUtils.Url
	DetectRst.OriginUrlStatusCode = respUtils.StatusCode
	DetectRst.OriginWebTitle = respUtils.WebTitle
	DetectRst.ContentLength = respUtils.ContentLength
	DetectRst.SiteUp = Up

	var redirectRespUtils p1httputils.HttpRespUtils
	//var _ *http.Response
	_, redirectUrl, isRedirect := p1httputils.CheckPageRedirect(resp, respUtils)
	if isRedirect {
		if !strings.Contains(redirectUrl, "http") {
			redirectUrl = fixedUrl + redirectUrl
		}
		_, redirectRespUtils, _ = p1httputils.HttpGet(redirectUrl, r.ProxyClient)

		//switch redirectType {
		//case "Location":
		//	// no need to check request err or not
		//	if !strings.Contains(redirectUrl, "http") {
		//		redirectUrl = fixedUrl + redirectUrl
		//	}
		//	_, redirectRespUtils, _ = p1httputils.HttpGet(redirectUrl, r.ProxyClient)
		//case "metaRefresh":
		//
		//case "jsRedirect":
		//	// todo
		//case "VueRoute":
		//	// todo
		//}

		DetectRst.RedirectUrl = redirectUrl
		DetectRst.RedirectWebTitle = redirectRespUtils.WebTitle
		P1fingerResps = append(P1fingerResps, redirectRespUtils)
	}

	for _, finger := range r.P1FingerPrints.GetElements() {
		for _, matcher := range finger.Matchers {
			matchFlag := false
			for _, targetInfo := range P1fingerResps {
				var content string
				switch matcher.Location {
				case "title":
					content = targetInfo.WebTitle
				case "header":
					if matcher.Type == "regex" {
						re := regexp.MustCompile(strings.Join(matcher.Words, "|")) // 预编译正则
						if re.MatchString(targetInfo.HeaderStr) {
							matchFlag = true
							finger.Name = re.FindString(targetInfo.HeaderStr)
						}
					} else {
						content = targetInfo.HeaderStr
					}
				case "body":
					content = targetInfo.BodyStr
				default:
					continue
				}

				if matcher.Type != "regex" && matchCondition(content, matcher.Words, matcher.Condition) {
					matchFlag = true
				}

				if matchFlag {
					DetectRst.FingerTag = append(DetectRst.FingerTag, finger.Name)
					break
				}
			}

		}
	}

	if len(DetectRst.FingerTag) <= 0 {
		// 规则匹配未匹配到，尝试主动路径匹配
		webPathMatched, DetectRstTmp := r.matchWithWebPath(fixedUrl)
		if webPathMatched {
			DetectRst = DetectRstTmp
		} else {
			r.DetectRstTdSafe.AddElement(DetectRst)
			return
		}
	}

	r.DetectRstTdSafe.AddElement(DetectRst)
	return
}

func (r *RuleClient) matchWithWebPath(fixedUrl string) (matchFlag bool, DetectRst DetectResult) {

	//var DetectRst DetectResult
	var p1RespUtils []p1httputils.HttpRespUtils

	for _, finger := range r.P1FingerPrints.GetElements() {
		for _, matcher := range finger.Matchers {
			if matcher.Location == "webPath" {
				fixWebPathUrl := fixedUrl + matcher.Path
				// 首次访问禁止重定向，手动解析重定向
				resp, rdrResputils, err := p1httputils.HttpGet(fixWebPathUrl, r.ProxyNoRedirectCilent)
				if err != nil {
					DetectRst = DetectResult{
						OriginUrl:           rdrResputils.Url,
						OriginUrlStatusCode: rdrResputils.StatusCode,
						OriginWebTitle:      rdrResputils.WebTitle,
						SiteUp:              Down,
						FingerTag:           []string{"WebSite down"},
					}
					r.DetectRstTdSafe.AddElement(DetectRst)
					return
				}
				p1RespUtils = append(p1RespUtils, rdrResputils)

				var P1fingerRedirectResp p1httputils.HttpRespUtils
				_, rdrUrl, isRedirect := p1httputils.CheckPageRedirect(resp, rdrResputils)
				if isRedirect {
					_, P1fingerRedirectResp, err = p1httputils.HttpGet(rdrUrl, r.ProxyClient)
					if err != nil {
						r.DetectRstTdSafe.AddElement(DetectRst)
						return
					}
					p1RespUtils = append(p1RespUtils, P1fingerRedirectResp)

					//switch redirectType {
					//case "Location":
					//
					//
					//case "jsRedirect":
					//	// todo
					//case "VueRoute":
					//	// todo
					//}
				}

				for _, fingerResp := range p1RespUtils {
					if matchCondition(fingerResp.BodyStr, matcher.Words, matcher.Condition) {
						DetectRst = DetectResult{
							OriginUrl:           rdrResputils.Url,
							OriginUrlStatusCode: rdrResputils.StatusCode,
							OriginWebTitle:      rdrResputils.WebTitle,
							SiteUp:              Up,
							FingerTag:           []string{finger.Name},
						}
						matchFlag = true
						return
					}
				}

			}
		}
	}

	return
}

func matchCondition(content string, words []string, condition string) bool {
	shooted := 0
	loweredContent := strings.ToLower(content)

	for _, word := range words {
		loweredWord := strings.ToLower(word)
		if strings.Contains(loweredContent, loweredWord) {
			shooted++
			if condition == "or" || condition == "" {
				return true
			}
		} else if condition == "and" {
			return false
		}
	}

	return condition == "and" && shooted == len(words)
}
