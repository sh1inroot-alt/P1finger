package cmd

import (
	"github.com/P001water/P1finger/cmd/vars"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	RootCmd.PersistentFlags().StringVarP(&vars.Options.ProxyUrl, "proxy", "p", "", "proxy eg: [--proxy socks5://127.0.0.1]")
	RootCmd.PersistentFlags().StringVarP(&vars.Options.Output, "output", "o", "p1finger.csv", "output file name: [-o p1finger.csv] / [-o p1finger.json]")
	RootCmd.PersistentFlags().BoolVar(&vars.Options.Debug, "debug", false, "http debug info, eg:[-debug]")
}

var RootCmd = &cobra.Command{
	Use:   "P1finger",
	Short: "一款红队行动下的重点资产指纹识别工具",
	Long: `
	██████╗  ██╗███████╗██╗███╗   ██╗ ██████╗ ███████╗██████╗ 
	██╔══██╗███║██╔════╝██║████╗  ██║██╔════╝ ██╔════╝██╔══██╗
	██████╔╝╚██║█████╗  ██║██╔██╗ ██║██║  ███╗█████╗  ██████╔╝
	██╔═══╝  ██║██╔══╝  ██║██║╚██╗██║██║   ██║██╔══╝  ██╔══██╗
	██║      ██║██║     ██║██║ ╚████║╚██████╔╝███████╗██║  ██║
	╚═╝      ╚═╝╚═╝     ╚═╝╚═╝  ╚═══╝ ╚═════╝ ╚══════╝╚═╝  ╚═
 
        一款红队行动下的重点资产指纹识别工具, Powered by P001water
`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// load P1finger config from p1fingerConf.yaml file
		filePath := filepath.Join(vars.ExecDir, "p1fingerConf.yaml")
		err := vars.LoadAppConf(filePath, &vars.AppConf)
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return
		}
	},
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:  RootCmd,
		Headings: cc.Red + cc.Bold + cc.Underline,
		Commands: cc.Cyan + cc.Bold,
		Example:  cc.Italic,
		ExecName: cc.Magenta + cc.Bold,
		Flags:    cc.Cyan + cc.Bold,
	})
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
