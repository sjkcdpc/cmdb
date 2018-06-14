package cmd

import (
	"github.com/mds1455975151/cmdb/server"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start cmdb server",
	Long:  "不想写了",
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
