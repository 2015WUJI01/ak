package progressbar

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
)

type ProgressBar = progressbar.ProgressBar

// New desc 行首描述，length 进度条格子长度，opts 官方支持的配置项
func New(desc string, length int, opts ...progressbar.Option) *progressbar.ProgressBar {
	defopts := []progressbar.Option{
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() { fmt.Println() }),
		progressbar.OptionSetWidth(25),
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
