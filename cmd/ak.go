package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use: "ak",
	// Use: "ak <command> [flags] <args>",
	// SilenceUsage: true,
	// Aliases:      nil,
	// SuggestFor:   nil,
	// Short: "",
	// Long:  "",
	// Example:      "ak opr 令",
	// Args:         cobra.MinimumNArgs(1),
	// Run: func(cmd *cobra.Command, args []string) {
	//
	// },
	// DisableFlagsInUseLine: true,
}

var (
	name string
	// alias string
	// id    int
)

func init() {
	rootCmd.AddCommand(updateCmd)
	// rootCmd.AddCommand(versionCmd)
	// rootCmd.SetUsageTemplate(RootCmdUsageTemplate())
	// rootCmd.SetHelpTemplate(RootCmdHelpTemplate())
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	// rootCmd.PersistentFlags().StringVarP(&name, "name", "n", "", "name usage")
	// rootCmd.LocalFlags().StringVarP(&alias, "alias", "a", "", "alias usage")
	// rootCmd.InheritedFlags().IntVarP(&id, "id", "i", 0, "id usage")
}
func Execute() error {
	return rootCmd.Execute()
}

//
// func RootCmdUsageTemplate() string {
// 	return HelperMsg()
// }
