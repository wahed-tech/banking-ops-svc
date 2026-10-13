package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wahed-tech/banking-ops-svc/conf"
	"github.com/wahed-tech/go-lib/v2/log"
)

var rootCmd = &cobra.Command{
	Use:   "banking-ops-svc",
	Short: "Banking ops back-office service",
	Long:  `Backs the ops back-office portal for banking: reconciliation, settlements and reporting APIs.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(conf.Init, initLogging)
}

func initLogging() {
	log.Init(conf.R.Logger.Level)
}
