package fofa

import (
	"github.com/P001water/P1finger/cmd"
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/P001water/P1finger/libs/fileutils"
	"github.com/P001water/P1finger/libs/sliceopt"
	"github.com/P001water/P1finger/modules/FofaClient"
	"github.com/P001water/P1finger/modules/RuleClient"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"strings"
)

var (
	Url     string
	UrlFile string
)

func init() {
	cmd.RootCmd.AddCommand(fofaCmd)
	fofaCmd.Flags().StringVarP(&vars.Options.Url, "url", "u", "", "target url")
	fofaCmd.Flags().StringVarP(&vars.Options.UrlFile, "file", "f", "", "target url file")
}

var fofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "Fingerprint Detect based on the Fofa cyberspace mapping engine.",
	Long:  "Fingerprint Detect based on the Fofa cyberspace mapping engine.",
	Run: func(cmd *cobra.Command, args []string) {
		gologger.Info().Msgf("p1fingeprint detect model: Fofa\n")
		err := fofaAction()
		if err != nil {
			gologger.Error().Msg(err.Error())
			return
		}
	},
}

func fofaAction() (err error) {
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
		gologger.Error().Msg("input url is null")
		return
	}
	targets = sliceopt.SliceRmDuplication(targets)

	vars.FofaCli, err = FofaClient.NewClient(
		FofaClient.WithURL("https://fofa.info/?email=&key=&version=v1"),
		FofaClient.WithAccountDebug(true),
		FofaClient.WithDebug(vars.Options.Debug),
		FofaClient.WithEmail(vars.AppConf.FofaCredentials.Email),
		FofaClient.WithApiKey(vars.AppConf.FofaCredentials.ApiKey),
	)

	// 美化查询语法用于Fofa
	var group []string
	var FinalQuery []string
	var querybeautify []string
	domains, ips := FofaClient.SplitDomainsAndIPs(targets)
	querybeautify = append(querybeautify, domains...)
	querybeautify = append(querybeautify, ips...)
	for i, simpleQuery := range querybeautify {
		group = append(group, simpleQuery)
		if (i+1)%50 == 0 || i == len(querybeautify)-1 {
			FinalQuery = append(FinalQuery, strings.Join(group, " || "))
			group = nil
		}
	}

	// 开始查询
	queryFields := []string{"ip", "port", "title", "product", "lastupdatetime", "protocol", "host"}
	for _, item := range FinalQuery {
		res, err := vars.FofaCli.HostSearch(item, -1, queryFields)
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return err
		}

		for _, simpleTarget := range res {
			tmp := RuleClient.DetectResult{
				OriginUrl:      simpleTarget[5] + "://" + simpleTarget[0] + ":" + simpleTarget[1],
				Host:           simpleTarget[6],
				OriginWebTitle: simpleTarget[2],
				FingerTag:      strings.Split(simpleTarget[3], ","),
				LastUpdateTime: simpleTarget[4],
			}
			vars.DetectResultTdSafe.AddElement(tmp)
		}
	}

	err = RuleClient.SaveToFile(vars.DetectResultTdSafe.GetElements(), vars.Options.Output)
	if err != nil {
		gologger.Error().Msg(err.Error())
		return
	}

	return nil
}
