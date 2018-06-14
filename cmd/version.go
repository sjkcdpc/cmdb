package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cmdb",
	Long:  `All software has versions, This is cmdb's'`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cmdb v1.0.0 -- develp(00a2ad5)")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
