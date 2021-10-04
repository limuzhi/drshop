package cmd

import (
	"drpshop/cmd/admin"
	"drpshop/cmd/sys"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "drpshop",
	Short:        "drpshop",
	SilenceUsage: true,
	Long:         `drpshop`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("可以使用")
	},
}

func init() {
	rootCmd.AddCommand(admin.StartCmd)
	rootCmd.AddCommand(sys.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
