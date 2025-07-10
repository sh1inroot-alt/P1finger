package vars

import (
	"github.com/P001water/P1finger/cmd/flag"
	"github.com/P001water/P1finger/modules/FofaClient"
	"github.com/P001water/P1finger/modules/RuleClient"
	"os"
	"path/filepath"
)

var (
	DetectResultTdSafe RuleClient.DetectResultTdSafeType //探测结果需要的数据
)

var (
	Options = &flag.Options{}

	ExecDir = filepath.Dir(os.Args[0])
	FofaCli *FofaClient.FofaClient
	AppConf P1fingerConf

	RstChan chan RuleClient.DetectResult
	Rsts    [][]string
)
