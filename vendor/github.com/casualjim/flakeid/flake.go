/*
github.com/boundary/flake in golang
*/

package flakeid

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"

	"regexp"

	"net"
	"strconv"
	"sync"
	"time"
)

type id struct {
	time uint64
	mac  uint64
	seq  uint64
}

func (i *id) Bytes() []byte {
	t := make([]byte, 8)
	s := make([]byte, 8)
	a := make([]byte, 16)
	binary.BigEndian.PutUint64(t, i.time)
	binary.BigEndian.PutUint64(a[6:14], i.mac)
	binary.BigEndian.PutUint64(s, i.seq)

	copy(a[0:6], t[2:8])
	copy(a[14:16], s[6:8])

	return a
}

func (i *id) Hex() string {
	return hex.EncodeToString(i.Bytes())
}

const (
	nano   = 1e6
	maxIds = 65365
)

const (
	sequenceBits   = 16
	workerIDShift  = 48
	timestampShift = 64
	sequenceMask   = -1 ^ (-1 << sequenceBits)
)

var (
	macStripRegexp = regexp.MustCompile(`[^a-fA-F0-9]`)
)

func (sf *Flake) generateHexIds(count int) ([]string, error) {
	ids, err := sf.generateIds(count)

	if err != nil {
		return nil, err
	}

	hexIds := make([]string, len(ids))

	for i := 0; i < count; i++ {
		hexIds[i] = ids[i].Hex()
	}

	return hexIds, nil
}

func (sf *Flake) generateIds(count int) ([]id, error) {
	if count < 1 {
		count = 1
	} else if count > 65365 {
		count = maxIds
	}

	ids := make([]id, count)

	sf.lock.Lock()
	defer sf.lock.Unlock()

	newTimeInMs := timestamp()

	if newTimeInMs > sf.lastTimestamp {
		sf.lastTimestamp = newTimeInMs
		sf.sequence = 0
	} else if newTimeInMs < sf.lastTimestamp {
		return nil, fmt.Errorf("Time has reversed! Old time: %v - New time: %v", sf.lastTimestamp, newTimeInMs)
	}

	for i := 0; i < count; i++ {
		sf.sequence++
		ids[i] = id{uint64(sf.lastTimestamp), sf.workerID, sf.sequence}
	}

	return ids, nil
}

// DefaultGenerator for flake ids
var DefaultGenerator = NewFlake()

// MustNewID creates a new ID or panics when a failure occurs
func MustNewID() string {
	id, err := NewID()
	if err != nil {
		panic(err)
	}
	return id
}

// NewID generates a new ID or an error
func NewID() (string, error) {
	return DefaultGenerator.Next()
}

// Flake represents a structure to build an id string from
type Flake struct {
	lastTimestamp int64
	workerID      uint64
	sequence      uint64
	lock          *sync.Mutex
}

// Next generates the next flake id
func (sf *Flake) Next() (string, error) {
	r, err := sf.generateHexIds(1)
	if err != nil {
		return "", err
	}
	return r[0], nil
}

// NextN generates the next n flake ids
func (sf *Flake) NextN(n int) ([]string, error) {
	return sf.generateHexIds(n)
}

// NewFlake creates a new instance of a flake id generator
func NewFlake() *Flake {
	return &Flake{workerID: DefaultWorkID(), lastTimestamp: -1, lock: new(sync.Mutex)}
}

func timestamp() int64 {
	return time.Now().Unix()
}

func tilNextMillis(ts int64) int64 {
	i := timestamp()
	for i <= ts {
		i = timestamp()
	}
	return i
}

// DefaultWorkID generates the id for this host
func DefaultWorkID() uint64 {
	ifs, err := net.Interfaces()

	if err != nil {
		log.Fatalf("Could not get any network interfaces: %v, %+v", err, ifs)
	}

	var hwAddr net.HardwareAddr

	for _, i := range ifs {
		if len(i.HardwareAddr) > 0 {
			hwAddr = i.HardwareAddr
			break
		}
	}

	if hwAddr == nil {
		log.Fatalf("No interface found with a MAC address: %+v", ifs)
	}

	mac := hwAddr.String()
	hex := macStripRegexp.ReplaceAllLiteralString(mac, "")

	u, err := strconv.ParseUint(hex, 16, 64)

	if err != nil {
		log.Fatalf("Unable to parse %v (from mac %v) as an integer: %v", hex, mac, err)
	}

	return u
}
