package persist

// Value returned from a persistence medium
type Value struct {
	Value       []byte
	Version     uint64
	LastUpdated int64
	_           struct{}
}
