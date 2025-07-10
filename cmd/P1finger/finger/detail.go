package finger

import (
	"fmt"
	"github.com/P001water/P1finger/cmd/vars"
	"github.com/P001water/P1finger/modules/RuleClient"
	"github.com/projectdiscovery/gologger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

func init() {
	fingersCmd.AddCommand(detailCmd)
}

var detailCmd = &cobra.Command{
	Use:   "detail",
	Short: "查看一个指纹的详细信息",
	Run: func(cmd *cobra.Command, args []string) {

		fingerName := args[0]

		fingerDetail(fingerName)

	},
}

func fingerDetail(fingerName string) {
	p1ruleClient, err := RuleClient.NewRuleClientBuilder().
		WithOutputFormat(vars.Options.Output).
		WithTimeout(10 * time.Second).
		Build()
	if err != nil {
		gologger.Error().Msgf("%v", err)
	}

	fingers := p1ruleClient.P1FingerPrints.GetElements()

	// 用于存储符合条件的 finger 对象
	var matchedFingers []RuleClient.FingerprintsType

	// 遍历 fingers 切片
	for _, finger := range fingers {
		if strings.Contains(finger.Name, fingerName) {
			matchedFingers = append(matchedFingers, finger)
		}
	}

	// 打开文件以写入
	file, err := os.Create("output.yaml")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// 将匹配的 finger 对象切片转换为 YAML 格式
	yamlData, err := yaml.Marshal(matchedFingers)
	if err != nil {
		fmt.Printf("Error converting to YAML: %v\n", err)
		return
	}

	// 写入文件
	_, err = file.Write(yamlData)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}

	fmt.Println("YAML data has been written to output.yaml")
}
