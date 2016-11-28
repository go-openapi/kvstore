package main

import (
	"log"
	"os"

	loads "github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/viper"

	"github.com/go-openapi/kvstore"
	"github.com/go-openapi/kvstore/api/handlers"
	"github.com/go-openapi/kvstore/gen/restapi"
	"github.com/go-openapi/kvstore/gen/restapi/operations"
)

func main() {
	cfg := viper.New()
	cfg.SetDefault("store.path", "./db/data.db")
	if err := cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalln(err)
		}
	}

	rt, err := kvstore.NewRuntime(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewKvstoreAPI(swaggerSpec)
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

	server.SetHandler(api.Serve(nil))

	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
