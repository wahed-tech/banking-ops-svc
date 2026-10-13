package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/wahed-tech/banking-ops-svc/conf"
	"github.com/wahed-tech/go-lib/v2/cmds"
	"github.com/wahed-tech/go-lib/v2/handler/commands"
	"github.com/wahed-tech/go-lib/v2/tracing"
)

var commandsWatcherCMD = &cobra.Command{
	Use:   "commands-watcher",
	Short: "Consume events related to commands",
	RunE:  cmds.WithErrorLog(consumeCommands),
}

func consumeCommands(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	handleSysCall(cancel)
	closer := tracing.InitTracer(conf.R.App.Name, conf.R.Tracing)
	defer closer.Close()

	commands.NewAccountHandler(conf.R.App.Name, conf.R.Kafka.ConfigMap).Start(ctx, commands.AccountHandlerFunc{
		OnRedactData:  nil, // Implement if needed
		OnEmailUpdate: nil, // Implement if needed
	})

	return nil
}
func handleSysCall(cancel context.CancelFunc) {
	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		<-gracefulStop
		cancel()
	}()
}
func init() {
	rootCmd.AddCommand(commandsWatcherCMD)
}
