package progressbar

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
)

type ProgressBar = progressbar.ProgressBar

func New(desc string, length int, opts ...progressbar.Option) *progressbar.ProgressBar {
	defopts := []progressbar.Option{
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
		progressbar.OptionSetWidth(50),
		// progressbar.OptionClearOnFinish(),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: "[green]-[reset]",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	}
	return progressbar.NewOptions(length, append(defopts, opts...)...)
}
