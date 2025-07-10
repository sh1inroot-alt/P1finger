package finger

import (
	"fmt"
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/P001water/P1finger/modules/RuleClient"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"time"
)

func init() {
	fingersCmd.AddCommand(listCmd)

}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有指纹的ID和Name",
	Run: func(cmd *cobra.Command, args []string) {
		p1ruleClient, err := RuleClient.NewRuleClientBuilder().
			WithRuleMode(vars.AppConf.RuleMode).
			WithOutputFormat(vars.Options.Output).
			WithTimeout(10 * time.Second).
			Build()
		if err != nil {
			gologger.Error().Msgf("%v", err)
		}

		fingers := p1ruleClient.P1FingerPrints.GetElements()

		for _, finger := range fingers {
			//f := fmt.Sprintf("%s--%s--%s", finger.ID, finger.Name, finger.FingerFile)
			f := fmt.Sprintf("%s", finger.Name)
			fmt.Println(f)
		}
	},
}
