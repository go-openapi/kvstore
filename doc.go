// Package kvstore represents a key/value database
package kvstore

//go:generate swagger generate client -A kvstore -t gen -f ./swagger/swagger.yml
//go:generate swagger generate server --exclude-main -A kvstore -t gen -f ./swagger/swagger.yml
//go:generate msgp -file ./persist/types.go
