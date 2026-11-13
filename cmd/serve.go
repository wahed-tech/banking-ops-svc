package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
	"github.com/wahed-tech/banking-ops-svc/api"
	"github.com/wahed-tech/banking-ops-svc/auth"
	"github.com/wahed-tech/banking-ops-svc/conf"
	"github.com/wahed-tech/banking-ops-svc/proc"
	"github.com/wahed-tech/banking-ops-svc/repository"
	"github.com/wahed-tech/banking-ops-svc/schema"
	"github.com/wahed-tech/go-lib/v2/cmds"
	"github.com/wahed-tech/go-lib/v2/handler"
	"github.com/wahed-tech/go-lib/v2/jwt"
	"github.com/wahed-tech/go-lib/v2/prom"
	"github.com/wahed-tech/go-lib/v2/server"
	"github.com/wahed-tech/go-lib/v2/skafka"
	"github.com/wahed-tech/go-lib/v2/tracing"
	"golang.org/x/sync/errgroup"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the REST server and Kafka consumer",
	Long:  `Starts the REST server used by the ops back-office portal and the Kafka consumer.`,
	RunE:  cmds.WithErrorLog(serve),
}

func serve(cmd *cobra.Command, _ []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	prom.Init()
	closer := tracing.InitTracer(conf.R.App.Name, conf.R.Tracing)
	defer closer.Close()

	pgPool, err := repository.NewPostgresPool(ctx, conf.R.Postgres)
	if err != nil {
		return err
	}
	defer pgPool.Close()

	keyFunc, err := jwt.NewKeyFuncURL([]string{conf.R.Auth.JWKSURL})
	if err != nil {
		return err
	}

	authenticator, err := auth.NewAuthenticator(keyFunc, conf.R.Auth)
	if err != nil {
		return err
	}

	gqlSchema, err := schema.New()
	if err != nil {
		return err
	}

	router := api.NewRouter(api.NewPingHandler(pgPool), authenticator, api.GraphQL(&gqlSchema), api.ParseOrigins(conf.R.Auth.AllowedOrigins))
	go server.Start(cancel, router, conf.R.Server.Host, conf.R.Server.Port)

	producer := skafka.NewProducer(conf.R.Kafka.ConfigMap)
	defer producer.Close()

	grp, ctx := errgroup.WithContext(ctx)
	startConsumer(ctx, grp, conf.R.Kafka.Topics.ExampleEvents, createHandlers(producer), producer)

	return grp.Wait()
}

func startConsumer(ctx context.Context, grp *errgroup.Group, topic string, handler *handler.Aggregate, producer skafka.Producer) {
	grp.Go(func() error {
		skafka.NewConsumer(conf.R.Kafka.ConfigMap, topic).StartAggregate(ctx, handler.Handle, skafka.RaiseErrorsAndTryReprocess(producer, conf.R.Kafka.Topics.ExampleEvents), time.Second*60)
		return nil
	})
}

func createHandlers(producer skafka.Producer) *handler.Aggregate {
	return handler.
		NewAggregate(handler.WithTracing(handler.WithLog(handler.SimpleStrategy()))).
		Register(
			handler.NewErrorAdapter(proc.NewCreateExample(producer)),
		)
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
