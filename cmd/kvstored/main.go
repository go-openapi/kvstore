package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	app "github.com/casualjim/go-app"
	"github.com/casualjim/middlewares"
	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/justinas/alice"

	"github.com/go-openapi/kvstore"
	"github.com/go-openapi/kvstore/api/handlers"
	"github.com/go-openapi/kvstore/gen/restapi"
	"github.com/go-openapi/kvstore/gen/restapi/operations"
)

func main() {

	app, err := app.New("kvstore")
	if err != nil {
		logrus.Fatalln(err)
	}

	log := app.Logger()
	cfg := app.Config()
	cfg.SetDefault("store.path", "./db/data.db")

	rt, err := kvstore.NewRuntime(app)
	if err != nil {
		log.Fatalln(err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewKvstoreAPI(swaggerSpec)
	api.Logger = log.Infof

	server := restapi.NewServer(api)
	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = `K/V store`
	parser.LongDescription = `K/V store is a simple single node store for retrieving key/value information`

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	api.KvDeleteEntryHandler = handlers.NewDeleteEntry(rt)
	api.KvFindKeysHandler = handlers.NewFindKeys(rt)
	api.KvGetEntryHandler = handlers.NewGetEntry(rt)
	api.KvPutEntryHandler = handlers.NewPutEntry(rt)

	handler := alice.New(
		middlewares.NewRecoveryMW(app.Info().Name, log),
		middlewares.NewAuditMW(app.Info(), log),
		middlewares.NewProfiler,
		middlewares.NewHealthChecksMW(app.Info().BasePath),
	).Then(api.Serve(nil))

	server.SetHandler(handler)

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
