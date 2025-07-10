package flag

type Options struct {
	Url     string
	UrlFile string

	ProxyUrl string // proxy address

	FingerDir  string
	FingerOnly string

	Debug bool
	Rate  int

	Update bool
	Output string
}
