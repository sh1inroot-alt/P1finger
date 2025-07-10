package p1print

import (
	"fmt"
	"github.com/P001water/P1finger/cmd/vars"
)

const (
	DebugPrefix = "[DBG] "
)

const (
	WEBTITLE       = "[WebInfo]"
	ERRTITLE       = "[ERRInfo]"
	UserInputTitle = " > "
)

// 包裹函数，添加前缀并调用原始的 fmt.Println
func UserInput(args ...interface{}) {
	fmt.Printf("\033[1;31;40m%v\033[0m ", UserInputTitle)
	fmt.Println(args...)
}

// 包裹函数，添加前缀并调用原始的 fmt.Printf
func UserInputf(format string, args ...interface{}) {
	fmt.Printf("\033[1;31;40m%v\033[0m ", UserInputTitle)
	fmt.Printf(format, args...)
}

// 包裹函数，添加前缀并调用原始的 fmt.Println
func ERRTitle(args ...interface{}) {
	fmt.Print(ERRTITLE)
	fmt.Println(args...)
}

// 包裹函数，添加前缀并调用原始的 fmt.Printf
func ERRTitlef(format string, args ...interface{}) {
	fmt.Print(ERRTITLE)
	fmt.Printf(format, args...)
}

// 包裹函数，添加前缀并调用原始的 fmt.Println
func WebInfo(args ...interface{}) {
	fmt.Print(WEBTITLE)
	fmt.Println(args...)
}

// 包裹函数，添加前缀并调用原始的 fmt.Printf
func WebInfof(format string, args ...interface{}) {
	fmt.Print(WEBTITLE)
	fmt.Printf(format, args...)
}

func Debug(args ...interface{}) {
	if vars.Options.Debug {
		fmt.Print(DebugPrefix)
		fmt.Println(args...)
	}
}

func Debugf(format string, args ...interface{}) {
	if vars.Options.Debug {
		fmt.Print(DebugPrefix)
		fmt.Printf(format, args...)
	}
}

//func main() {
//	// 示例用法
//	Info("This is information.")
//	Fail("This is failure.")
//	Warn("This is warning.")
//	Error("This is an error.")
//
//	Infof("This is formatted info: %s\n", "Hello, World!")
//	Errorf("This is formatted error: %d\n", 404)
//}
