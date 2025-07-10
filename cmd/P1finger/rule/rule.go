package rule

import (
	"fmt"
	"github.com/P001water/P1finger/libs/p1print"
	"github.com/P001water/P1finger/libs/sliceopt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/P001water/P1finger/cmd"
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/P001water/P1finger/libs/fileutils"
	"github.com/P001water/P1finger/modules/RuleClient"
	"github.com/k0kubun/go-ansi"
	"github.com/projectdiscovery/gologger"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

func init() {
	cmd.RootCmd.AddCommand(RuleCmd)

	RuleCmd.Flags().StringVarP(&vars.Options.Url, "url", "u", "", "target url")
	RuleCmd.Flags().StringVarP(&vars.Options.UrlFile, "file", "f", "", "target url file")
	RuleCmd.Flags().IntVar(&vars.Options.Rate, "rate", 500, "The number of go coroutines")
}

var RuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Fingerprint Detect based on the P1finger local fingerprint database",
	Run: func(cmd *cobra.Command, args []string) {
		err := RuleRun()
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
	},
}

// 按照 OriginUrlStatusCode 对 DetectResult 切片进行排序
type ByStatusCode []RuleClient.DetectResult

func (a ByStatusCode) Len() int           { return len(a) }
func (a ByStatusCode) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStatusCode) Less(i, j int) bool { return a[i].OriginUrlStatusCode < a[j].OriginUrlStatusCode }

func RuleRun() (err error) {

	p1ruleClient, err := RuleClient.NewRuleClientBuilder().
		WithRuleMode(vars.AppConf.RuleMode).
		WithProxyURL(vars.Options.ProxyUrl).
		WithOutputFormat(vars.Options.Output).
		WithTimeout(10 * time.Second).
		Build()
	if err != nil {
		return
	}

	// 整合目标输入
	var targets []string
	if vars.Options.Url != "" {
		targets = append(targets, vars.Options.Url)
	}

	if vars.Options.UrlFile != "" {
		var urlsFromFile []string
		filePath := vars.Options.UrlFile
		urlsFromFile, err = fileutils.ReadLinesFromFile(filePath)
		if err != nil {
			return
		}
		targets = append(targets, urlsFromFile...)
	}

	if len(targets) <= 0 {
		return fmt.Errorf("targets can't be null")
	}
	targets = sliceopt.SliceRmDuplication(targets)

	Progressbar := progressbar.NewOptions(len(targets),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[cyan][P1finger][reset] Detection progress..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	concurrency := vars.Options.Rate
	var workWg sync.WaitGroup
	urlChan := make(chan string, len(targets))
	vars.RstChan = make(chan RuleClient.DetectResult, len(targets))

	// 生产 - 发送任务到 channel
	for _, url := range targets {
		workWg.Add(1)
		urlChan <- url
	}

	results := []RuleClient.DetectResult{}
	go func() {
		for r := range vars.RstChan {
			results = append(results, r)
			workWg.Done()
		}
	}()

	// 消费-启动固定数量的 worker
	for i := 0; i < concurrency; i++ {
		go func() {
			for consumerUrl := range urlChan {
				defer workWg.Done()
				dtctRst, err := p1ruleClient.Detect(consumerUrl)
				if err != nil {
					dtctRst.OriginUrlStatusCode = 999
					workWg.Add(1)
					vars.RstChan <- dtctRst
					return
				}
				workWg.Add(1)
				vars.RstChan <- dtctRst
				Progressbar.Add(1)
			}
		}()
	}

	close(urlChan)
	workWg.Wait()
	Progressbar.Finish()
	fmt.Println()

	sort.Sort(ByStatusCode(results))

	Reset := "\033[0m"  // 重置颜色
	Green := "\033[32m" // 绿色
	Red := "\033[31m"   // 绿色
	for _, r := range results {
		statusColor := statusCodeColor(r.OriginUrlStatusCode)
		fingerTags := strings.Join(r.FingerTag, ",")
		if len(r.RedirectUrl) > 0 {
			fmt.Printf("%s %s [%s] Redirect to %v [%v] %s%s%s\n", statusColor, r.OriginUrl, r.OriginWebTitle, r.RedirectUrl, r.RedirectWebTitle, Green, fingerTags, Reset)
		} else if r.OriginUrlStatusCode == 999 {
			fmt.Printf("%s %s Error:%s%s%s\n", statusColor, r.OriginUrl, Red, fingerTags, Reset)
		} else if len(r.FingerTag) > 0 {
			fmt.Printf("%s %s [%s] %s%s%s\n", statusColor, r.OriginUrl, r.OriginWebTitle, Green, fingerTags, Reset)
		} else {
			fmt.Printf("%s %s [%s]\n", statusColor, r.OriginUrl, r.OriginWebTitle)
		}
	}

	p1print.Debugf("从文件导入目标总数：%v\n", len(targets))
	p1print.Debugf("文件导入去重后目标总数：%v\n", len(sliceopt.SliceRmDuplication(targets)))
	p1print.Debugf("管道已经检查目标总数：%v\n", len(p1ruleClient.DetectRstTdSafe.GetElements()))

	// save to file
	err = RuleClient.SaveToFile(results, p1ruleClient.OutputFormat)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	return

}

// 根据状态码返回颜色化的字符串
func statusCodeColor(statusCode int) string {
	switch {
	case statusCode >= 200 && statusCode < 300:
		return "\033[32m" + fmt.Sprintf("[%d]", statusCode) + "\033[0m" // 绿色
	case statusCode >= 300 && statusCode < 400:
		return "\033[34m" + fmt.Sprintf("[%d]", statusCode) + "\033[0m" // 蓝色
	case statusCode >= 400 && statusCode < 500:
		return "\033[33m" + fmt.Sprintf("[%d]", statusCode) + "\033[0m" // 黄色
	case statusCode >= 500:
		return "\033[31m" + fmt.Sprintf("[%d]", statusCode) + "\033[0m" // 红色
	default:
		return fmt.Sprintf("[%d]", statusCode)
	}
}
