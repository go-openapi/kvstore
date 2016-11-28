package handlers

import (
	"io/ioutil"

	"github.com/go-openapi/kvstore"
	"github.com/go-openapi/kvstore/gen/restapi/operations/kv"
	"github.com/go-openapi/kvstore/persist"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// NewPutEntry handles a request for saving an entry
func NewPutEntry(rt *kvstore.Runtime) kv.PutEntryHandler {
	return &putEntry{rt: rt}
}

type putEntry struct {
	rt *kvstore.Runtime
}

// Handle the put entry request
func (d *putEntry) Handle(params kv.PutEntryParams) middleware.Responder {
	rid := swag.StringValue(params.XRequestID)
	key := params.Key

	value, err := ioutil.ReadAll(params.Body)
	e := params.Body.Close()
	if err != nil {
		return kv.NewPutEntryDefault(0).WithXRequestID(rid).WithPayload(modelsError(err))
	}
	if e != nil {
		return kv.NewPutEntryDefault(0).WithXRequestID(rid).WithPayload(modelsError(e))
	}

	if err := d.rt.DB().Put(key, value); err != nil {
		if err == persist.ErrNotFound {
			return kv.NewPutEntryNotFound().WithXRequestID(rid).WithPayload(modelsError(err))
		}
		return kv.NewPutEntryDefault(0).WithXRequestID(rid).WithPayload(modelsError(err))
	}
	return kv.NewPutEntryNoContent().WithXRequestID(rid)
}
