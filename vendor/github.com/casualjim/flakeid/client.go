package flakeid

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
)

var newLine = []byte("\n")

// Client for a flake id server
type Client struct {
	io.Closer
	conn net.Conn
}

// NewClient creates a new client connectio to a flake id server
func NewClient(proto, addr string, verbose bool) (*Client, error) {
	if verbose {
		log.Printf("connecting to %s://%s", proto, addr)
	}
	conn, err := net.Dial(proto, addr)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// Close connection with server
func (c *Client) Close() error {
	return c.conn.Close()
}

// NextN gets the next n ids from a server
func (c *Client) NextN(n uint32) ([]string, error) {
	data := make([]byte, 4)
	binary.BigEndian.PutUint32(data, n)

	if _, err := c.conn.Write(data); err != nil {
		return nil, err
	}

	result := make([]byte, int(n)*33)
	if _, err := c.conn.Read(result); err != nil {
		return nil, err
	}

	res := make([]string, 0, int(n))
	for _, ln := range bytes.SplitN(result, newLine, -1) {
		res = append(res, string(ln))
	}
	return res, nil
}
