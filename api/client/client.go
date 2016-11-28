package client

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	httpclient "github.com/go-openapi/kvstore/gen/client"
	"github.com/go-openapi/kvstore/gen/client/kv"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/swag"
)

// New client for the kv store api
func New(uri string) (*KvStore, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	return &KvStore{client: httpclient.New(httptransport.New(u.Host, u.Path, []string{u.Scheme}), nil)}, nil
}

// KvStore wraps the swagger client for central handling of error cases etc
type KvStore struct {
	client *httpclient.Kvstore
}

// FindKeys for a given prefix
func (k *KvStore) FindKeys(prefix string) ([]string, error) {
	var pref *string
	if prefix != "" {
		pref = swag.String(prefix)
	}
	keys, err := k.client.Kv.FindKeys(kv.NewFindKeysParams().WithPrefix(pref))
	if err != nil {
		return nil, err
	}
	return keys.Payload, nil
}

// Delete an entry from the store
func (k *KvStore) Delete(key string) error {
	_, err := k.client.Kv.DeleteEntry(kv.NewDeleteEntryParams().WithKey(key))
	if err != nil {
		return fmt.Errorf("failed to delete %q because: %v", key, err)
	}
	return nil
}

// Entry in the k/v store
type Entry struct {
	// Data the payload to save
	Data []byte
	// Version is required when this is an update
	Version uint64
	_       struct{}
}

// Put an entry in the k/v store
func (k *KvStore) Put(key string, data *Entry) error {
	params := kv.NewPutEntryParams().WithKey(key).WithBody(bytes.NewBuffer(data.Data))
	if data.Version != 0 {
		params.SetIfMatch(swag.String(strconv.FormatUint(data.Version, 10)))
	}

	created, updated, err := k.client.Kv.PutEntry(params)
	if err != nil {
		switch e := err.(type) {
		case *kv.PutEntryConflict:
			return errors.New(swag.StringValue(e.Payload.Message))
		case *kv.PutEntryNotFound:
			return errors.New(swag.StringValue(e.Payload.Message))
		case *kv.PutEntryGone:
			return errors.New(swag.StringValue(e.Payload.Message))
		case *kv.PutEntryDefault:
			return errors.New(swag.StringValue(e.Payload.Message))
		default:
			return e
		}
	}

	var etag string
	if updated != nil {
		etag = updated.ETag
	} else if created != nil {
		etag = created.Etag
	}

	if etag != "" {
		v, err := strconv.ParseUint(etag, 10, 64)
		if err != nil {
			return err
		}
		data.Version = v
	}
	return nil
}

// Get a value from the store, when the version not 0 it will use that to
// get a not modified response.
func (k *KvStore) Get(key string, version uint64) (*Entry, error) {
	params := kv.NewGetEntryParams().WithKey(key)
	if version != 0 {
		params.SetIfNoneMatch(swag.String(strconv.FormatUint(version, 10)))
	}

	data := bytes.NewBuffer(nil)
	value, err := k.client.Kv.GetEntry(params, data)
	if err != nil {
		switch e := err.(type) {
		case *kv.GetEntryNotFound:
			return nil, errors.New(swag.StringValue(e.Payload.Message))
		case *kv.GetEntryNotModified:
			return &Entry{Version: version}, nil
		case *kv.GetEntryDefault:
			return nil, errors.New(swag.StringValue(e.Payload.Message))
		default:
			return nil, e
		}
	}

	entry := new(Entry)
	entry.Data = data.Bytes()
	if value.ETag != "" {
		v, err := strconv.ParseUint(value.ETag, 10, 64)
		if err != nil {
			return nil, err
		}
		entry.Version = v
	}
	return entry, nil
}
