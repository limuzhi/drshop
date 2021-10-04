package main

import (
	"github.com/spf13/cobra"
	"log"
	"whole/pkg/kratos-cli/internal/proto"
)

var (
	version = "v0.0.1"

	rootCmd = &cobra.Command{
		Use:     "kratos-cli",
		Short:   "kratos 开发扩展工具",
		Long:    `Kratos 开发扩展工具`,
		Version: version,
	}
)

func init() {
	rootCmd.AddCommand(proto.CmdProto)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
