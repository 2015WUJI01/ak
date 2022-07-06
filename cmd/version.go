package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const version = "0.0.1"

func init() {
	// rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:    "version",
	Short:  "输出 ak 当前的版本信息",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("ak version %s\n", version)
		return nil
	},
	DisableFlagsInUseLine: true,
}
