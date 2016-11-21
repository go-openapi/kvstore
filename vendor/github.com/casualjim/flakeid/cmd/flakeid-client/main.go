package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/casualjim/flakeid"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	Address    string `long:"address" short:"a" description:"the address to connect to" default:"[::1]:3525"`
	SocketPath string `long:"socket-path" short:"s" description:"the unix domain socket to listen on" default:"/var/run/flakeid.sock"`
	Unix       bool   `long:"unix" description:"when present will listen on unix domain socket"`
	Count      uint32 `long:"count" short:"c" description:"amount of ids to generate" default:"1"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Fatalln(err)
	}

	var client *flakeid.Client
	var err error
	if opts.Unix {
		client, err = flakeid.NewClient("unix", opts.SocketPath, false)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		client, err = flakeid.NewClient("tcp", opts.Address, false)
		if err != nil {
			log.Fatalln(err)
		}
	}

	result, err := client.NextN(opts.Count)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(strings.Join(result, "\n"))

	if err := client.Close(); err != nil {
		log.Fatalln(err)
	}
}
