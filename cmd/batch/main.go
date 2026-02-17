package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"sb-scanner/cmd/batch/cmd"
)

var rootCmd = &cobra.Command{
	Use:   "batch",
	Short: "SB Scanner Batch",
	Long:  "SB Scanner Batch Service",
}

func main() {
	flags := rootCmd.Flags()
	flags.String("config", "", "path to config file")
	rootCmd.PersistentFlags().AddFlagSet(flags)
	viper.BindPFlags(flags)

	rootCmd.AddCommand(cmd.Sync())
	rootCmd.Execute()
}
