package client

import (
	"net/url"

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

func (k *KvStore) Delete(key string) error {

	k.client.Kv.DeleteEntry(kv.NewDeleteEntryParams().WithKey(key))
}
