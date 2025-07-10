package p1print

import (
	"fmt"
)

func IfAddColor(text string, NoColor bool) (colorText string) {
	if NoColor {
		colorText = fmt.Sprintf("%v", text)
		//showStr += InfoCMDData
	} else {
		colorText = fmt.Sprintf("\u001B[1;31;40m%v\u001B[0m", text)
	}
	return
}
