package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/casualjim/flakeid"
	"github.com/jessevdk/go-flags"
)

const (
	unix    = "unix"
	tcp     = "tcp"
	newLine = '\n'
)

var opts struct {
	Port       int    `long:"port" short:"p" env:"PORT" description:"the port to listen on, defaults to 3525" default:"3525"`
	SocketPath string `long:"socket-path" short:"s" description:"the unix domain socket to listen on" default:"/var/run/flakeid.sock"`
	NoTCP      bool   `long:"no-tcp" description:"when specified no tcp listener will be started"`
	Unix       bool   `long:"unix" description:"when present will listen on unix domain socket"`
	// HTTP       bool   `long:"http" description:"when present will listen on HTTP"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Fatalln(err)
	}

	if opts.Unix {
		listener, err := net.Listen(unix, opts.SocketPath)
		defer listener.Close()
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("serving flake ids at unix://%s\n", listener.Addr())
		go listen(listener)
	}

	if !opts.NoTCP {
		listener, err := net.Listen(tcp, ":"+strconv.Itoa(opts.Port))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("serving flake ids at tcp://%s\n", listener.Addr())
		listen(listener)
	}
}

func listen(l net.Listener) {
	for {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatalln(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	// read from transport
	data := make([]byte, 4)
	_, err := conn.Read(data)
	if err != nil {
		log.Fatalln(err)
	}

	count := binary.BigEndian.Uint32(data)
	if err != nil {
		log.Fatalln(err)
	}

	// do the actual work
	ids, err := flakeid.DefaultGenerator.NextN(int(count))
	if err != nil {
		log.Fatalln(err)
	}

	// marshal
	buf := bytes.NewBuffer(make([]byte, 0, len(ids)*33-1))
	for i, id := range ids {
		if i > 0 {
			buf.WriteRune(newLine)
		}
		buf.WriteString(id)
	}

	// reply
	if _, err := conn.Write(buf.Bytes()); err != nil {
		log.Fatalln(err)
	}
}
