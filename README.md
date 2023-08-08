# K/V store

Is a very simple open API based rest API to provide a key value store

## Generation

To generate the code required this application makes use of:

* [go-swagger](https://goswagger.io)
* [msgpack generator](https://github.com/tinylib/msgp)

Install with:

```
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
go install github.com/tinylib/msgp@latest
```

All the generation happens based on `go generate` and is configured in `$project_root/doc.go`.

This application chooses to configure itself differently than the default generated code because it's more convenient and allows for a nicer structure.
The configuration happens in `$project_root/cmd/kvstored`.

## Install

You can install this application with the regular go means.

```
go install github.com/go-openapi/kvstore/cmd/...@latest
```

## Running

```
$ ./kvstored --help
Usage:
  kvstored [OPTIONS]

K/V store is a simple single node store for retrieving key/value information

Application Options:
      --scheme=            the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec
      --cleanup-timeout=   grace period for which to wait before shutting down the server (default: 10s)
      --max-header-size=   controls the maximum number of bytes the server will read parsing the request header's keys and values, including
                           the request line. It does not limit the size of the request body. (default: 1MiB)
      --socket-path=       the unix socket to listen on (default: /var/run/kvstore.sock)
      --host=              the IP to listen on (default: localhost) [$HOST]
      --port=              the port to listen on for insecure connections, defaults to a random value [$PORT]
      --listen-limit=      limit the number of outstanding requests
      --keep-alive=        sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop
                           mid-download) (default: 3m)
      --read-timeout=      maximum duration before timing out read of the request (default: 30s)
      --write-timeout=     maximum duration before timing out write of the response (default: 60s)
      --tls-host=          the IP to listen on for tls, when not specified it's the same as --host [$TLS_HOST]
      --tls-port=          the port to listen on for secure connections, defaults to a random value [$TLS_PORT]
      --tls-certificate=   the certificate to use for secure connections [$TLS_CERTIFICATE]
      --tls-key=           the private key to use for secure conections [$TLS_PRIVATE_KEY]
      --tls-listen-limit=  limit the number of outstanding requests
      --tls-keep-alive=    sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections ( e.g. closing laptop
                           mid-download)
      --tls-read-timeout=  maximum duration before timing out read of the request
      --tls-write-timeout= maximum duration before timing out write of the response

Help Options:
  -h, --help               Show this help message
```
